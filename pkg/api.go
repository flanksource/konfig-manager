package pkg

import (
	"fmt"
	"net/http"

	"github.com/flanksource/commons/logger"
	"k8s.io/apimachinery/pkg/util/json"
)

type APIServer struct {
	Repos      []string
	Branches   []string
	ConfigFile string
}

func GetConfigData(input *APIServer, appName string, showObjects string) ([]map[string]KustomizeResources, error) {
	hierarchy, err := GetHierarchy(input.ConfigFile, appName)
	if err != nil {
		return nil, err
	}
	data := ReadConfiguration(input.Repos, input.Branches, hierarchy)

	if showObjects == "" {
		data = removeObjectsFromList(data)
	}

	return data, nil
}

func (input *APIServer) GetConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		showObjects := req.URL.Query().Get("objects")
		appName := req.URL.Query().Get("application")

		data, err := GetConfigData(input, appName, showObjects)
		if err != nil {
			logger.Fatalf("failed to get hierarchy: %v", err)
		}

		output, err := json.Marshal(data)
		if err != nil {
			logger.Errorf("error marshalling the data: %w", err)
			return
		}
		_, err = w.Write(output)
		if err != nil {
			fmt.Println(err)
			logger.Errorf("error writing the body: %w", err)
			return
		}
	}
}

func removeObjectsFromList(dataWithKustomizeObjects []map[string]KustomizeResources) []map[string]KustomizeResources {
	var data []map[string]KustomizeResources
	for _, objMap := range dataWithKustomizeObjects {
		for key, kustomizeResource := range objMap {
			kustomizeResource.Objects = nil
			objMap[key] = kustomizeResource
		}
		data = append(data, objMap)
	}
	return data
}
