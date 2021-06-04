package pkg

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"gopkg.in/yaml.v3"
)

func getObjectsFromHierarchy(config Config, resources []Resource) []Resource {
	var objs []Resource
	for _, item := range config.Hierarchy {
		resource := getResourceFromHierarchy(resources, item)
		if resource == nil {
			continue
		}
		resource.Hieratchy = item
		objs = append(objs, *resource)
	}
	return objs
}

func getResourceFromHierarchy(resources []Resource, item Item) *Resource {
	for _, resource := range resources {
		if resource.Item.GetKind() == item.Kind &&
			resource.Item.GetName() == item.Name &&
			resource.Item.GetNamespace() == item.Namespace {
			return &resource
		}
	}
	return nil
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
