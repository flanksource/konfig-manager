package pkg

import (
	"net/http"

	"github.com/flanksource/commons/logger"
	"k8s.io/apimachinery/pkg/util/json"
)

type ServerInput struct {
	Repos     []string
	Branches  []string
	Hierarchy Config
}

func (input *ServerInput) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	showObjects := req.URL.Query().Get("objects")
	data := GenerateMetrics(input.Repos, input.Branches, input.Hierarchy)
	if showObjects == "" {
		data = removeObjectsFromList(data)
	}
	output, err := json.Marshal(data)
	if err != nil {
		logger.Fatalf("error marshalling the data")
		return
	}
	_, err = w.Write(output)
	if err != nil {
		logger.Fatalf("Error writing the body")
		return
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
