package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/flanksource/konfig-manager/pkg"
	"github.com/magiconair/properties"
)

type testInput struct {
	config        string            // location of config file with hierarchy rules
	data          string            // location of the data file with the manifests
	applications  []string          // name of application to be generated
	verifications map[string]string // Map of all the properties need to be checked with the given value
}

var tests = map[string]testInput{
	"testHierarchyMergeWithStdinInput": {
		config: "fixtures/spring-config.yml",
		data:   "fixtures/spring.yml",
		verifications: map[string]string{
			"config-key":                  "some-value",
			"config-key-quotes":           "some-other-value",
			"spring.datasource.maxActive": "40",
			"null-string-key":             "null",
			"undefined-string-key":        "undefined",
			"bool-string-key":             "true",
			"int-string-key":              "11",
		},
		applications: []string{"spring"},
	},
	"testHierarchyMergeWithInputFile": {
		config: "fixtures/spring-config.yml",
		data:   "fixtures/spring.yml",
		verifications: map[string]string{
			"config-key":                  "some-value",
			"spring.datasource.maxActive": "40",
		},
		applications: []string{"spring"},
	},
	"testReadFromConfigMapCreatedWithFile": {
		config: "fixtures/fileProperties-config.yml",
		data:   "fixtures/fileProperties.yml",
		verifications: map[string]string{
			"some-key":                              "value-from-spring",
			"new-key":                               "diff-value",
			"logging.level.org.springframework.web": "INFO",
		},
		applications: []string{"spring"},
	},
	"testSecretValues": {
		config: "fixtures/secret-config.yml",
		data:   "fixtures/data-with-secrets.yml",
		verifications: map[string]string{
			"secret-key":                            "some-value",
			"logging.level.org.springframework.web": "INFO",
		},
		applications: []string{"spring"},
	},
	"testMultipleApplications": {
		config: "fixtures/multi-application-config.yml",
		data:   "fixtures/multi-applications.yml",
		verifications: map[string]string{ // putting common configs here which would be present in both the properties file
			"config-key": "some-value",
			"secret-key": "some-value",
			"new-key":    "diff-value",
		},
		applications: []string{"spring", "quarkus"},
	},
}

func TestGenerateJs(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resources, err := pkg.ReadResources(test.data)
			if err != nil {
				t.Error(err)
			}

			// Adding an empty resource tests the ability of the function to exclude them without panicking.
			resources = append([]pkg.Resource{pkg.Resource{}}, resources...)

			for _, name := range test.applications {
				hierarchy, err := pkg.GetHierarchy(test.config, name)
				if err != nil {
					t.Error(err)
				}

				file := hierarchy.GenerateJsPropertiesFile(resources)

				p := properties.MustLoadString(file)
				// check property key and values

				for key, value := range test.verifications {
					var transformedValue string
					if _, err := strconv.Atoi(value); err == nil {
						transformedValue = fmt.Sprintf("%v;", value)
					} else if _, err := strconv.ParseBool(value); err == nil {
						transformedValue = fmt.Sprintf("%v;", value)
					} else if value == "null" || value == "undefined" {
						transformedValue = fmt.Sprintf("%v;", value)
					} else {
						transformedValue = fmt.Sprintf("\"%v\";", value)
					}

					transformedKey := fmt.Sprintf("window['__%v__']", key)
					propVal, exists := p.Get(transformedKey)
					if !exists {
						t.Errorf("property not found: %s", key)
					}
					if propVal != transformedValue {
						t.Errorf("%s: expected %s got %s", key, propVal, transformedValue)
					}
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resources, err := pkg.ReadResources(test.data)
			if err != nil {
				t.Error(err)
			}

			// Adding an empty resource tests the ability of the function to exclude them without panicking.
			resources = append([]pkg.Resource{pkg.Resource{}}, resources...)

			for _, name := range test.applications {
				hierarchy, err := pkg.GetHierarchy(test.config, name)
				if err != nil {
					t.Error(err)
				}
				file := hierarchy.GeneratePropertiesFile(resources)

				p := properties.MustLoadString(file)
				// check property key and values
				for key, value := range test.verifications {
					propVal, exists := p.Get(key)
					if !exists {
						t.Errorf("property not found: %s", key)
					}
					if propVal != value {
						t.Errorf("%s: expected %s got %s", key, propVal, value)
					}
				}
			}
		})
	}
}
