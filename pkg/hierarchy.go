package pkg

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"

	"github.com/pkg/errors"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/commons/text"
	"github.com/flanksource/kommons"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Hierarchy    []Item   `yaml:"hierarchy" json:"hierarchy"`
	Applications []string `yaml:"applications,omitempty" json:"applications,omitempty"`
}

type Item struct {
	Kind          string       `yaml:"kind" json:"kind"`
	Name          string       `yaml:"name" json:"name"`
	Namespace     string       `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Type          ResourceType `yaml:"type,omitempty" json:"type,omitempty"`
	Key           string       `yaml:"key,omitempty" json:"key,omitempty"`
	Index         int          `yaml:"index,omitempty" json:"index,omitempty"`
	HierarchyName string       `yaml:"hierarchyName,omitempty" json:"hierarchyName,omitempty"`
}

func (item Item) String() string {
	s := fmt.Sprintf("%s/%s", item.Kind, item.Name)
	if item.Key != "" {
		s += fmt.Sprintf("[%s]", item.Key)
	}
	return s
}

func (item Item) FindIn(resources []Resource) *Resource {
	for _, resource := range resources {
		if resource.Item != nil && resource.Item.GetKind() == item.Kind &&
			resource.Item.GetName() == item.Name &&
			(item.Namespace == "" || item.Namespace == resource.Item.GetNamespace()) {
			return &resource
		}
	}
	return nil
}

func (config Config) GetPropertiesMap(resources []Resource) map[string]string {
	propertiesMap := make(map[string]string)
	for _, property := range config.GetProperties(resources) {
		propertiesMap[property.Key] = property.Value
	}
	return propertiesMap
}

func (config Config) GetProperties(resources []Resource) map[string]Property {
	propertiesMap := make(map[string]Property)
	for _, resource := range config.WalkHierarchy(resources) {
		for _, v := range resource.GetProperties() {
			propertiesMap[v.Key] = v
		}
	}
	return propertiesMap
}

func (config Config) getPropertiesBySection(resources []Resource) map[Item]Properties {
	var bySection = make(map[Item]Properties)
	for _, v := range config.GetProperties(resources) {
		if _, ok := bySection[v.Resource.Hierarchy]; !ok {
			bySection[v.Resource.Hierarchy] = Properties{}
		}
		bySection[v.Resource.Hierarchy] = append(bySection[v.Resource.Hierarchy], v)
	}
	return bySection
}

func (config Config) GeneratePropertiesFile(resources []Resource) string {
	var properties string
	var bySection = config.getPropertiesBySection(resources)

	for _, item := range config.Hierarchy {
		list := bySection[item]
		sort.Sort(list)
		if len(list) == 0 {
			continue
		}
		properties += fmt.Sprintf("#\n# %s\n#\n", item.String())
		for _, property := range list {
			properties += fmt.Sprintf("%v=%v\n", property.Key, property.Value)
		}
	}

	return properties
}

func (config Config) GenerateJsPropertiesFile(resources []Resource) string {
	var properties string
	var bySection = config.getPropertiesBySection(resources)
	for _, item := range config.Hierarchy {
		list := bySection[item]
		sort.Sort(list)
		if len(list) == 0 {
			continue
		}
		properties += fmt.Sprintf("//\n// %s\n//\n", item.String())
		for _, property := range list {
			if _, err := strconv.Atoi(property.Value); err == nil {
				properties += fmt.Sprintf("window['__%v__']=%v;\n", property.Key, property.Value)
			} else if _, err := strconv.ParseBool(property.Value); err == nil {
				properties += fmt.Sprintf("window['__%v__']=%v;\n", property.Key, property.Value)
			} else if property.Value == "null" || property.Value == "undefined" {
				properties += fmt.Sprintf("window['__%v__']=%v;\n", property.Key, property.Value)
			} else {
				properties += fmt.Sprintf("window['__%v__']=\"%v\";\n", property.Key, property.Value)
			}
		}
	}

	return properties
}

func (config Config) WalkHierarchy(resources []Resource) []Resource {
	var objs []Resource
	for _, item := range config.Hierarchy {
		logger.Tracef("[%s] finding in %d resources", item, len(resources))
		resource := item.FindIn(resources)
		if resource == nil {
			continue
		}
		logger.Infof("[%s] found %s", item, kommons.GetName(resource.Item))
		resource.Hierarchy = item
		objs = append(objs, *resource)
	}
	return objs
}

func GetHierarchy(configFile, applicationName string) (Config, error) {
	logger.Infof("[%s] getting hierarchy for %s", configFile, applicationName)
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, errors.Wrap(err, fmt.Sprintf("error reading %s", configFile))
	}
	var config Config
	data, err := text.Template(string(buf), map[string]string{"name": applicationName})
	logger.Tracef(data)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return Config{}, err
	}
	for i := range config.Hierarchy {
		config.Hierarchy[i].Index = i
	}
	return config, nil
}
