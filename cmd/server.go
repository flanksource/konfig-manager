package cmd

import (
	"github.com/flanksource/konfig-manager/pkg"
	"github.com/spf13/cobra"
)

var Server = &cobra.Command{
	Use:   "server",
	Short: "Start the konfig manager server and exposes the metrics",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Server(cmd)
	},
}

func init() {
	Server.Flags().StringP("config-file", "c", "konfig-manager.yaml", "specify the config file")
	Server.Flags().Int("port", 8080, "http port")
	Server.Flags().StringSliceP("repos", "r", []string{}, "list of repos to parse")
	Server.Flags().StringSliceP("branches", "b", []string{"main"}, "list of branches to parse in the specified repos")
	Server.Flags().StringP("allowed-origins", "", "", "To set the allowed origins in the http server")
}
