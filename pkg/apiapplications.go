package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/flanksource/commons/logger"
	"gopkg.in/yaml.v3"
)

func (input *APIServer) GetApplicationHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		handleError := func(errMessage string, err error, code int) {
			logger.Errorf("%s: %v", errMessage, err)
			newErr, _ := json.Marshal(&Error{ErrorMessage: errMessage})
			resp.WriteHeader(code)
			if _, err := resp.Write(newErr); err != nil {
				logger.Errorf("failed to write body: %v", err)
			}
		}

		applications, err := getApplicationNames(input.ConfigFile)
		if err != nil {
			handleError("failed to get application names", err, http.StatusNotImplemented)
		}
		output, err := json.Marshal(applications)
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

func getApplicationNames(configFile string) ([]string, error) {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}
	return config.Applications, nil
}
