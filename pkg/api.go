package pkg

import (
	"net/http"

	"github.com/flanksource/commons/logger"
	"k8s.io/apimachinery/pkg/util/json"
)

type APIServer struct {
	Repos      []string
	Branches   []string
	ConfigFile string
}

type Error struct {
	ErrorMessage string `yaml:"errorMessage" json:"errorMessage"`
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
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")

		showObjects := req.URL.Query().Get("objects")
		appName := req.URL.Query().Get("application")

		data, err := GetConfigData(input, appName, showObjects)

		handleError := func(errMessage string, err error, code int) {
			logger.Errorf("%s: %v", errMessage, err)
			newErr, _ := json.Marshal(&Error{ErrorMessage: errMessage})
			resp.WriteHeader(code)
			if _, err := resp.Write(newErr); err != nil {
				logger.Errorf("failed to write body: %v", err)
			}
		}

		if err != nil {
			handleError("failed to get hierarchy", err, http.StatusNotImplemented)
			return
		}

		output, err := json.Marshal(data)
		if err != nil {
			handleError("failed to marshal data", err, http.StatusNotImplemented)
			return
		}

		if _, err := resp.Write(output); err != nil {
			logger.Errorf("failed to write body: %v", err)
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
