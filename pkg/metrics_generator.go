package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flanksource/commons/exec"
	"github.com/flanksource/commons/logger"
	"github.com/flanksource/kommons/kustomize"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/rand"
)

type KustomizeResources struct {
	Resources   []string `yaml:"resources"`
	Bases       []string `yaml:"bases"`
	Environment string
	FilePath    string
	FileDir     string
	Branch      string
	SubResource bool
	ImportedBy  []string
	Global      bool
	Objects     []Resource
	Properties  map[string]string
}

func GenerateMetrics(repos []string, branches []string, hierarchy Config) []map[string]KustomizeResources {
	var data []map[string]KustomizeResources
	for _, repo := range repos {
		data = append(data, parseRepoWithBranches(repo, branches, hierarchy))
	}
	return data
}

func parseRepoWithBranches(repo string, branches []string, hierarchy Config) map[string]KustomizeResources {
	var data = make(map[string]KustomizeResources)
	dir, _ := ioutil.TempDir("/tmp", rand.String(5))
	err := exec.Exec(fmt.Sprintf("git clone %v %v", repo, dir))
	if err != nil {
		logger.Fatalf("error cloning repo %v in temp dir %v", repo, dir)
		return nil
	}
	logger.Debugf("successfully cloned the repo %v", repo)
	if err := os.Chdir(dir); err != nil {
		logger.Fatalf("error changing context dir %v", dir)
		return nil
	}
	logger.Debugf("changed context dir to %v", dir)
	for _, branch := range branches {
		var resourceList []KustomizeResources
		if err := exec.Exec(fmt.Sprintf("git checkout %v", branch)); err != nil {
			logger.Fatalf("error checking out branch %v", branch)
			return nil
		}
		if err := filepath.Walk(".",
			func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.Name() != "kustomization.yaml" {
					return nil
				}
				var resource = KustomizeResources{
					FilePath:    filePath, // this is absolute from the root
					FileDir:     path.Dir(filePath),
					Branch:      branch,
					Environment: filepath.Base(path.Dir(filePath)),
				}
				buf, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(buf, &resource); err != nil {
					return err
				}
				// adding objects only if present in the config
				resource.Objects = getObjectsFromHierarchy(hierarchy, getKustomizeResourceObjects(resource.Resources, resource.FileDir))
				data[resource.FileDir] = resource
				resourceList = append(resourceList, resource)
				return nil
			}); err != nil {
			logger.Fatalf("Error walking the current directory")
			return nil
		}
	}
	for _, resource := range data {
		for _, base := range resource.Bases {
			baseAbsolutePath := filepath.Join(resource.FileDir, base)
			other := data[baseAbsolutePath]
			if strings.HasPrefix(baseAbsolutePath, resource.FileDir) {
				other.SubResource = true
			} else {
				other.ImportedBy = append(data[baseAbsolutePath].ImportedBy, resource.FilePath)
				other.Global = true
			}
			data[baseAbsolutePath] = other
		}
	}

	data = fillProperties(data, hierarchy)
	return data
}

func fillProperties(dataWithOutProperties map[string]KustomizeResources, hierarchy Config) map[string]KustomizeResources {
	var dataWithProperties = make(map[string]KustomizeResources)
	for key, kustomizeResource := range dataWithOutProperties {
		kustomizeResource.Properties = generatePropertiesForApp(hierarchy, kustomizeResource.Objects)
		dataWithProperties[key] = kustomizeResource
	}
	return dataWithProperties
}

func getKustomizeResourceObjects(resourcesName []string, fileDir string) []Resource {
	var resources []Resource
	for _, resourceName := range resourcesName {
		filePath := filepath.Join(fileDir, resourceName)
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Errorf("error reading file %v: %v", filePath, err)
		}
		objs, err := kustomize.GetUnstructuredObjects(buf)
		if err != nil {
			logger.Errorf("error parsing objects: %v", err)
		}
		var resource []Resource
		for _, obj := range objs {
			resource = append(resource, Resource{
				Item: obj.(*unstructured.Unstructured),
				Path: filePath,
			})
		}
		resources = append(resources, resource...)
	}
	return resources
}
