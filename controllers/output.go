package controllers

import (
	konfigmanagerv1 "github.com/flanksource/konfig-manager/api/v1"
	"github.com/flanksource/konfig-manager/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	configMapKind  = "ConfigMap"
	secretKind     = "Secret"
	coreAPIVersion = "v1"
)

func (r *KonfigReconciler) createOutputObject(output konfigmanagerv1.Output, config pkg.Config, resources []pkg.Resource) error {
	properties := make(map[string]string)
	if output.Type == "file" || output.Type == "File" {
		if output.Key == "" {
			output.Key = "application.properties"
		}
		properties[output.Key] = config.GeneratePropertiesFile(resources)
	} else {
		properties = config.GetPropertiesMap(resources)
	}
	if output.Kind == "ConfigMap" || output.Kind == "configmap" || output.Kind == "cm" {
		if err := r.Kommons.Apply(output.Namespace, getConfigMap(output.Name, output.Namespace, properties)); err != nil {
			r.Log.Error(err, "error creating/updating configmap", output.Name, output.Namespace)
			return err
		}
		r.Log.Info("created/updated configmap", output.Name, output.Namespace)
		return nil
	}
	if output.Kind == "Secret" || output.Kind == "secret" {
		propertiesWithBytes := make(map[string][]byte)
		for key, value := range properties {
			propertiesWithBytes[key] = []byte(value)
		}
		if err := r.Kommons.Apply(output.Namespace, getSecret(output.Name, output.Namespace, propertiesWithBytes)); err != nil {
			r.Log.Error(err, "error creating/updating secret", output.Name, output.Namespace)
			return err
		}
		r.Log.Info("created/updated secret", output.Name, output.Namespace)
	}
	return nil
}

func getConfigMap(name, namespace string, properties map[string]string) runtime.Object {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       configMapKind,
			APIVersion: coreAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: properties,
	}
}

func getSecret(name, namespace string, properties map[string][]byte) runtime.Object {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       secretKind,
			APIVersion: coreAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: properties,
	}
}

func (r *KonfigReconciler) getResources(config pkg.Config) ([]pkg.Resource, error) {
	var resources []pkg.Resource
	for _, item := range config.Hierarchy {
		obj, err := r.Kommons.GetByKind(item.Kind, item.Namespace, item.Name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, pkg.Resource{Item: obj})
	}
	return resources, nil
}
