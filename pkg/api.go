package pkg

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

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

func (input *APIServer) GetSpaHandler(rootFs http.FileSystem) http.HandlerFunc {
	readFile := func(file http.File) ([]byte, error) {
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		fileBuf := make([]byte, size)
		_, err := file.Read(fileBuf)
		return fileBuf, err
	}

	return func(resp http.ResponseWriter, req *http.Request) {
		path, err := filepath.Abs(req.URL.Path)
		if err != nil {
			// if we failed to get the absolute path respond with a 400 bad request
			// and stop
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = rootFs.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			// if the path does not return a file or directory, serve back index.html
			indexFile, err := rootFs.Open("/index.html")
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			fileData, err := readFile(indexFile)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := resp.Write(fileData); err != nil {
				logger.Errorf("failed to write body: %v", err)
			}
			return
		} else if err != nil {
			// if we got an error (that wasn't that the file doesn't exist) stating the
			// file, return a 500 internal server error and stop
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		http.FileServer(rootFs).ServeHTTP(resp, req)
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
