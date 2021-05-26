/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"

	"github.com/flanksource/konfig-manager/cmd"
	"github.com/spf13/cobra"
)

var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "konfig-manager",
	Short: "konfig-manager is responsible for managing configs based on hierarchy provided in the config file",
}

func main() {
	rootCmd.AddCommand(cmd.GenerateCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of konfig-manager",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})
	rootCmd.SetUsageTemplate(rootCmd.UsageTemplate() + fmt.Sprintf("\nversion: %s\n ", version))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
