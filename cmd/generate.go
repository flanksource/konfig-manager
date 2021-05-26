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

package cmd

import (
	"io/ioutil"

	"github.com/flanksource/konfig-manager/pkg"
	"github.com/spf13/cobra"
)

var (
	input, output, config string
	applicationNames      []string
)

// GenerateCmd represents the base command when called without any subcommands
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates application-properties based on the hierarchy and config provided",
	RunE: func(cmd *cobra.Command, args []string) error {
		var buf []byte
		var err error
		if input == "-" {
			buf, err = ioutil.ReadFile("/dev/stdin")
			if err != nil {
				return err
			}
		} else {
			buf, err = ioutil.ReadFile(input)
			if err != nil {
				return err
			}
		}
		if err := pkg.GenerateProperties(buf, applicationNames, config, output); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	GenerateCmd.PersistentFlags().StringVarP(&input, "input", "i", "-", "input of yaml dump")
	GenerateCmd.PersistentFlags().StringVarP(&config, "config", "c", "config.yml", "path to config file consisting hierarchy. Defaults to config.yml in pwd")
	GenerateCmd.PersistentFlags().StringSliceVarP(&applicationNames, "app-name", "A", []string{}, "name of application being templated")
	GenerateCmd.PersistentFlags().StringVarP(&output, "output", "o", ".", "path to directory where the properties file would be created. Defaults to pwd")
}
