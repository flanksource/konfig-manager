package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/flanksource/commons/console"
	"github.com/flanksource/commons/exec"
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
)

var (
	konfigManager = "bin/konfig-manager"
	outputDir     = "test/output/"
)
var testResults = &console.TestResults{
	Writer: os.Stdout,
}

var tests = map[string]Test{
	"testHierarchyMergeWithStdinInput":     TestHierarchyMergeWithStdinInput,
	"testHierarchyMergeWithInputFile":      TestHierarchyMergeWithInputFile,
	"testReadFromConfigMapCreatedWithFile": TestReadFromConfigMapCreatedWithFile,
	"testSecretValues":                     TestSecretValues,
	"testMultipleApplications":             TestMultipleApplications,
}

var inputs = map[string]testInput{
	"testHierarchyMergeWithStdinInput": {
		config:       "test/spring-config.yml",
		data:         "test/spring.yml",
		applications: "spring",
		verifications: map[string]string{
			"config-key":                  "some-value",
			"spring.datasource.maxActive": "40",
		},
	},
	"testHierarchyMergeWithInputFile": {
		config:       "test/spring-config.yml",
		data:         "test/spring.yml",
		applications: "spring",
		verifications: map[string]string{
			"config-key":                  "some-value",
			"spring.datasource.maxActive": "40",
		},
	},
	"testReadFromConfigMapCreatedWithFile": {
		config:       "test/fileProperties-config.yml",
		data:         "test/fileProperties.yml",
		applications: "spring",
		verifications: map[string]string{
			"some-key":                              "value-from-spring",
			"new-key":                               "diff-value",
			"logging.level.org.springframework.web": "INFO",
		},
	},
	"testSecretValues": {
		config:       "test/secret-config.yml",
		data:         "test/data-with-secrets.yml",
		applications: "spring",
		verifications: map[string]string{
			"secret-key":                            "some-value",
			"logging.level.org.springframework.web": "INFO",
		},
	},
	"testMultipleApplications": {
		config:       "test/multi-application-config.yml",
		data:         "test/multi-applications.yml",
		applications: "spring,quarkus",
		verifications: map[string]string{ // putting common configs here which would be present in both the properties file
			"config-key": "some-value",
			"secret-key": "some-value",
			"new-key":    "diff-value",
		},
	},
}

type Test func(*console.TestResults, testInput) error
type testInput struct {
	config        string            // location of config file with hierarchy rules
	data          string            // location of the data file with the manifests
	applications  string            // name of application to be generated
	verifications map[string]string // Map of all the properties need to be checked with the given value
}

func main() {
	errors := map[string]error{}

	for name, test := range tests {
		if err := test(testResults, inputs[name]); err != nil {
			errors[name] = err
		}
	}
	if len(errors) > 0 {
		for name, err := range errors {
			log.Errorf("test %s failed: %v", name, err)
		}
		os.Exit(1)
	}
	log.Infof("All tests passed !!!")
}

func TestHierarchyMergeWithInputFile(test *console.TestResults, input testInput) error {
	config := "test/spring-config.yml"
	data := "test/spring.yml"
	_ = os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)
	name := "testHierarchyMergeWithInputFile"
	application := "spring"
	command := fmt.Sprintf("%v generate -c %v -i %v -A %v -o %v", konfigManager, config, data, application, outputDir)
	if err := exec.Exec(command); err != nil {
		return err
	}

	// will verify if file exists
	err := verifyProperties(test, name, outputDir+input.applications+".properties", input.verifications)
	if err != nil {
		return err
	}
	test.Passf(name, "Test passed: override confirmed")
	return nil
}

func TestHierarchyMergeWithStdinInput(test *console.TestResults, input testInput) error {
	_ = os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	name := "testHierarchyMergeFromStdin"
	command := fmt.Sprintf("cat %v | %v generate -c %v -A %v -o %v", input.data, konfigManager, input.config, input.applications, outputDir)
	if err := exec.Exec(command); err != nil {
		return err
	}
	// will verify if file exists
	err := verifyProperties(test, name, outputDir+input.applications+".properties", input.verifications)
	if err != nil {
		return err
	}
	test.Passf(name, "Test passed: override confirmed")
	return nil
}

func TestReadFromConfigMapCreatedWithFile(test *console.TestResults, input testInput) error {
	_ = os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	name := "testReadFromConfigMapCreatedWithFile"
	command := fmt.Sprintf("cat %v | %v generate -c %v -A %v -o %v", input.data, konfigManager, input.config, input.applications, outputDir)
	if err := exec.Exec(command); err != nil {
		return err
	}
	// will verify if file exists
	err := verifyProperties(test, name, outputDir+input.applications+".properties", input.verifications)
	if err != nil {
		return err
	}
	test.Passf(name, "Test passed: config values read successfully")
	return nil
}

func TestSecretValues(test *console.TestResults, input testInput) error {
	_ = os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	name := "testSecretValues"
	command := fmt.Sprintf("cat %v | %v generate -c %v -A %v -o %v", input.data, konfigManager, input.config, input.applications, outputDir)
	if err := exec.Exec(command); err != nil {
		return err
	}
	// will verify if file exists
	err := verifyProperties(test, name, outputDir+input.applications+".properties", input.verifications)
	if err != nil {
		return err
	}
	test.Passf(name, "Test passed: secret values read successfully")
	return nil
}

func TestMultipleApplications(test *console.TestResults, input testInput) error {
	_ = os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	name := "testMultipleApplication"
	command := fmt.Sprintf("cat %v | %v generate -c %v -A %v -o %v", input.data, konfigManager, input.config, input.applications, outputDir)
	if err := exec.Exec(command); err != nil {
		return err
	}
	apps := strings.Split(input.applications, ",")

	// verify for applications
	for _, app := range apps {
		if err := verifyProperties(test, name, outputDir+app+".properties", input.verifications); err != nil {
			return err
		}
	}
	test.Passf(name, "Test passed: both property file exists with all the common configs")
	return nil
}

func verifyProperties(test *console.TestResults, testName, file string, verifications map[string]string) error {
	p := properties.MustLoadFile(file, properties.UTF8)
	// check property key and values
	for key, value := range verifications {
		propVal, exists := p.Get(key)
		if !exists {
			test.Failf(testName, "the key provided in the verifications input does not exist :%v", key)
			return fmt.Errorf("")
		}
		if propVal != value {
			test.Failf(testName, "The expected value for %v was %v but got %v", key, value, propVal)
			return fmt.Errorf("")
		}
	}
	return nil
}
