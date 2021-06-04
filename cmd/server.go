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
	Server.Flags().StringP("config-file", "c", "", "specify the config file")
	Server.Flags().Int("http-port", 8080, "port to expose the health dashboard")
	Server.Flags().StringSliceP("repos", "r", []string{}, "list of repos to parse")
	Server.Flags().StringSliceP("branches", "b", []string{}, "list of branches to parse in the specified repos")
	Server.Flags().StringP("allowedOrigins", "o", "", "To set the allowed origins in the http server")
}
