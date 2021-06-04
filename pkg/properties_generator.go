package pkg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/flanksource/kommons/kustomize"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func generatePropertiesForApp(config Config, resources []Resource) map[string]string {
	objects := getObjectsFromHierarchy(config, resources)
	propertiesMap := make(map[string]string)
	propertiesMapList := make([]map[string]string, 50)
	for _, resource := range objects {
		propertiesMapList = append(propertiesMapList, resource.GetPropertiesMap())
	}
	for _, propMap := range propertiesMapList {
		for key, value := range propMap {
			propertiesMap[key] = value
		}
	}
	return propertiesMap
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

func createPropetiesFile(propertiesMap map[string]string, output, applicationName string) error {
	var properties string
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

func GenerateProperties(buf []byte, applicationNames []string, config, output string) error {
	resources, err := getResources(buf)
	if err != nil {
		return err
	}
	for _, name := range applicationNames {
		hierarchy, err := getHierarchy(config, name)
		if err != nil {
			return err
		}
		propertiesMap := generatePropertiesForApp(hierarchy, resources)
		err = createPropetiesFile(propertiesMap, output, name)
		if err != nil {
			return err
		}
	}
	return nil
}
