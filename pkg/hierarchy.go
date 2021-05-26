package pkg

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"gopkg.in/yaml.v3"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getObjectsFromHierarchy(input Input) map[*unstructured.Unstructured]string {
	var objs = make(map[*unstructured.Unstructured]string)
	for _, item := range input.Config.Hierarchy {
		key, value := getResourceFromHierarchy(input.Resources, item)
		objs[key] = value
	}
	return objs
}

func getResourceFromHierarchy(resources []Resource, item Item) (*unstructured.Unstructured, string) {
	for _, resource := range resources {
		if resource.Item.GetKind() == item.Kind &&
			resource.Item.GetName() == item.Name &&
			resource.Item.GetNamespace() == item.Namespace {
			if item.Type == "file" || item.Type == "File" {
				return resource.Item, item.Key
			}
			return resource.Item, "data"
		}
	}

	return nil, ""
}

func getHierarchy(configFile, applicationName string) (Config, error) {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	var config Config
	data, err := templateHierarchy(applicationName, string(buf))
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func templateHierarchy(name, data string) ([]byte, error) {
	t, err := template.New("hierarchy").Parse(data)
	if err != nil {
		return []byte(""), err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, name)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}
