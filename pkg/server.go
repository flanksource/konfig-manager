package pkg

import (
	"fmt"
	nethttp "net/http"
	"path/filepath"

	"github.com/flanksource/commons/logger"
	"github.com/spf13/cobra"
)

func Server(cmd *cobra.Command) {
	configFilePath, _ := cmd.Flags().GetString("config-file")
	httpPort, _ := cmd.Flags().GetInt("port")
	repos, _ := cmd.Flags().GetStringSlice("repos")
	branches, _ := cmd.Flags().GetStringSlice("branches")
	allowedOrigins, _ := cmd.Flags().GetString("allowed-origins")

	configFilePathAbsPath, err := filepath.Abs(configFilePath)
	if err != nil {
		logger.Fatalf("failed to parse config file path: %v", err)
	}

	server := &APIServer{
		Repos:      repos,
		ConfigFile: configFilePathAbsPath,
		Branches:   branches,
	}
	handler := server.GetConfigHandler()
	applicationHandler := server.GetApplicationHandler()
	nethttp.HandleFunc("/api", simpleCors(handler, allowedOrigins))
	nethttp.HandleFunc("/api/applications", simpleCors(applicationHandler, allowedOrigins))
	addr := fmt.Sprintf("0.0.0.0:%d", httpPort)
	if err := nethttp.ListenAndServe(addr, nil); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}

// simpleCors is minimal middleware for injecting an Access-Control-Allow-Origin header value.
// If an empty allowedOrigin is specified, then no header is added.
func simpleCors(f nethttp.HandlerFunc, allowedOrigin string) nethttp.HandlerFunc {
	// if not set return a no-op middleware
	if allowedOrigin == "" {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {
			f(w, r)
		}
	}
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		f(w, r)
	}
}
