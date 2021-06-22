package test

import (
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
			"spring.datasource.maxActive": "40",
		},
	},
	"testHierarchyMergeWithInputFile": {
		config: "fixtures/spring-config.yml",
		data:   "fixtures/spring.yml",
		verifications: map[string]string{
			"config-key":                  "some-value",
			"spring.datasource.maxActive": "40",
		},
	},
	"testReadFromConfigMapCreatedWithFile": {
		config: "fixtures/fileProperties-config.yml",
		data:   "fixtures/fileProperties.yml",
		verifications: map[string]string{
			"some-key":                              "value-from-spring",
			"new-key":                               "diff-value",
			"logging.level.org.springframework.web": "INFO",
		},
	},
	"testSecretValues": {
		config: "fixtures/secret-config.yml",
		data:   "fixtures/data-with-secrets.yml",
		verifications: map[string]string{
			"secret-key":                            "some-value",
			"logging.level.org.springframework.web": "INFO",
		},
	},
	"testMultipleApplications": {
		config: "fixtures/multi-application-config.yml",
		data:   "fixtures/multi-applications.yml",
		verifications: map[string]string{ // putting common configs here which would be present in both the properties file
			"config-key": "some-value",
			"secret-key": "some-value",
			"new-key":    "diff-value",
		},
	},
}

func TestGenerate(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resources, err := pkg.ReadResources(test.data)
			if err != nil {
				t.Error(err)
			}

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
