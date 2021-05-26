package pkg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/flanksource/kommons/kustomize"
	"github.com/hairyhenderson/gomplate/v3/base64"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func generatePropertiesForApp(input Input, output, applicationName string) error {
	itemsMap := getObjectsFromHierarchy(input)
	propertiesMap := make(map[string]interface{})
	var properties string
	for item, key := range itemsMap {
		if key != "data" {
			properties = properties + item.Object["data"].(map[string]interface{})[key].(string)
		} else {
			data, bool, err := unstructured.NestedMap(item.Object, "data")
			if !bool || err != nil {
				return err
			}
			if item.GetKind() == "Secret" {
				for prop, value := range data {
					val, _ := base64.Decode(value.(string))
					propertiesMap[prop] = string(val)
				}
			} else {
				for prop, value := range data {
					propertiesMap[prop] = value
				}
			}
		}
	}
	for prop, value := range propertiesMap {
		properties = properties + fmt.Sprintf("%v=%v\n", prop, value)
	}
	var filePath string
	if strings.HasSuffix(output, "/") {
		filePath = output + applicationName + ".properties"
	} else {
		filePath = output + "/" + applicationName + ".properties"
	}

	err := ioutil.WriteFile(filePath, []byte(properties), 0755)
	if err != nil {
		return err
	}
	return nil
}

func getResources(buf []byte) ([]Resource, error) {
	resources, err := kustomize.GetUnstructuredObjects(buf)
	if err != nil {
		return nil, err
	}
	var inputResources []Resource
	for _, resource := range resources {
		inputResources = append(inputResources, Resource{
			Item: resource.(*unstructured.Unstructured),
		})
	}
	return inputResources, nil
}

func GenerateProperties(buf []byte, applicationNames []string, config, output string) error {
	resources, err := getResources(buf)
	if err != nil {
		return err
	}
	for _, name := range applicationNames {
		hierarchyItems, err := getHierarchyItems(config, name)
		if err != nil {
			return err
		}
		inputStruct := Input{
			Resources: resources,
			Hierarchy: Hierarchy{
				Items: hierarchyItems,
			},
			Applications: applicationNames,
		}
		err = generatePropertiesForApp(inputStruct, output, name)
		if err != nil {
			return err
		}
	}
	return nil
}
