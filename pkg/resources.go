package pkg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/flanksource/commons/logger"

	"github.com/flanksource/kommons/kustomize"
	"github.com/hairyhenderson/gomplate/v3/base64"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Resource struct {
	// path is the relative path to the file containing the resource definition
	Path string
	// A configmap, secret or sealed secret
	Item *unstructured.Unstructured
	// A type of resource that resource holds either native or properties from file
	Hierarchy Item
}

func (r Resource) String() string {
	return fmt.Sprintf("%s/%s/%s", r.Item.GetKind(), r.Item.GetNamespace(), r.Item.GetName())
}

type ResourceType string

const (
	ResourceTypeNative     ResourceType = "native"
	ResourceTypeProperties ResourceType = "properties"
)

type Property struct {
	Key, Value string
	Comment    string
	Resource   Resource
}

type Properties []Property

func (p Property) String() string {
	return fmt.Sprintf("%s=%s", p.Key, p.Value)
}
func (p Properties) Len() int {
	return len(p)
}

func (p Properties) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Properties) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}

func (r Resource) GetProperties() []Property {
	var properties []Property
	for k, v := range r.GetPropertiesMap() {
		properties = append(properties, Property{
			Key:      k,
			Value:    v,
			Comment:  r.String(),
			Resource: r,
		})
	}

	return properties
}

func (r Resource) GetPropertiesMap() map[string]string {
	if r.Hierarchy.Type == ResourceTypeProperties {
		return r.GeneratePropertesMapFromProperties()
	}
	return r.GeneratePropertiesMapFromNative()
}

func (r Resource) GeneratePropertesMapFromProperties() map[string]string {
	var propertiesMap = make(map[string]string)
	var prop string
	value := r.Item.Object["data"].(map[string]interface{})[r.Hierarchy.Key]
	if value == nil {
		logger.Debugf("can not find the given key: %v in the %v: %v", r.Hierarchy.Key, r.Item.GetKind(), r.Item.GetName())
		return nil
	}
	if r.Item.GetKind() == "Secret" {
		val, _ := base64.Decode(value.(string))
		prop = string(val)
	} else {
		prop = value.(string)
	}
	for _, keyValue := range strings.Split(prop, "\n") {
		propKeyValue := strings.Split(keyValue, "=")
		if len(propKeyValue) == 2 {
			propKey := propKeyValue[0]
			propValue := propKeyValue[1]
			propertiesMap[propKey] = propValue
		}
	}
	return propertiesMap
}

func (r Resource) GeneratePropertiesMapFromNative() map[string]string {
	var propertiesMap = make(map[string]string)
	data, ok, err := unstructured.NestedMap(r.Item.Object, "data")
	if !ok || err != nil {
		return nil
	}
	if r.Item.GetKind() == "Secret" {
		for prop, value := range data {
			val, _ := base64.Decode(value.(string))
			propertiesMap[prop] = string(val)
		}
	} else {
		for prop, value := range data {
			propertiesMap[prop] = value.(string)
		}
	}
	return propertiesMap
}

func ReadResources(input string) ([]Resource, error) {
	var buf []byte
	var err error
	if input == "-" {
		buf, err = ioutil.ReadFile("/dev/stdin")
		if err != nil {
			return nil, err
		}
	} else {
		buf, err = ioutil.ReadFile(input)
		if err != nil {
			return nil, err
		}
	}
	return GetResources(buf)
}

func GetResources(buf []byte) ([]Resource, error) {
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
