package controllers

import (
	"fmt"

	konfigmanagerv1 "github.com/flanksource/konfig-manager/api/v1"
	"github.com/flanksource/konfig-manager/pkg"
)

func (r *HierarchyConfigReconciler) createOutput(output konfigmanagerv1.Output, properties map[string]string) error {
	if output.Type == "file" || output.Type == "File" {
		properties = getPropertiesInFileFormat(properties, output.FileName)
	}
	if output.Kind == "ConfigMap" || output.Kind == "configmap" || output.Kind == "cm" {
		if err := r.Kommons.CreateOrUpdateConfigMap(output.Name, output.Namespace, properties); err != nil {
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
		if err := r.Kommons.CreateOrUpdateSecret(output.Name, output.Namespace, propertiesWithBytes); err != nil {
			return err
		}
		r.Log.Info("created/updated secret", output.Name, output.Namespace)
	}
	//r.Kommons.Apply()
	return nil
}

func getPropertiesInFileFormat(properties map[string]string, filename string) map[string]string {
	propertiesFile := make(map[string]string)
	if filename == "" {
		filename = "application.properties"
	}
	var data string
	for key, value := range properties {
		data = data + fmt.Sprintf("%v=%v\n", key, value)
	}
	propertiesFile[filename] = data
	return propertiesFile
}

func (r *HierarchyConfigReconciler) getResources(config pkg.Config) ([]pkg.Resource, error) {
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
