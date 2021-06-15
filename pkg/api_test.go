package pkg

import (
	"path/filepath"
	"testing"

	"github.com/flanksource/commons/logger"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigData(t *testing.T) {
	configFilePathAbsPath, err := filepath.Abs("../test/fixtures/large-config.yaml")
	if err != nil {
		logger.Fatalf("failed to parse config file path: %v", err)
	}

	type args struct {
		input       *APIServer
		appName     string
		showObjects string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		verifications map[string]map[string]string
	}{
		{
			name: "TestApplicationOne",
			args: args{
				input: &APIServer{
					Repos:      []string{"https://github.com/flanksource/testing-sample-repo.git"},
					Branches:   []string{"main"},
					ConfigFile: configFilePathAbsPath,
				},
				appName:     "one",
				showObjects: "true",
			},
			wantErr: false,
			verifications: map[string]map[string]string{
				"namespaces/foo-lala/apps/spam": {
					"URL_GHOST":             "https://ghost/omen/toast/",
					"a.datasource.username": "IUGIYGIU9",
				},
			},
		},
		{
			name: "TestApplicationTwo",
			args: args{
				input: &APIServer{
					Repos:      []string{"https://github.com/flanksource/testing-sample-repo.git"},
					Branches:   []string{"main"},
					ConfigFile: configFilePathAbsPath,
				},
				appName:     "two",
				showObjects: "true",
			},
			wantErr: false,
			verifications: map[string]map[string]string{
				"namespaces/foo-lala/apps/spam": {
					"NO_PROXY":       "",
					"URL_GHOST_DATA": "https://ghost:8080/",
					"URL_TOAST_HOST": "http://toast:8080/burnt",
					"no_proxy":       "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfigData(tt.args.input, tt.args.appName, tt.args.showObjects)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for path, checks := range tt.verifications {
				for check, value := range checks {
					assert.Equal(t, got[0][path].Properties[check], value)
				}
			}
		})
	}
}
