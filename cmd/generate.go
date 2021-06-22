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
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/commons/text"
	"github.com/flanksource/kommons"
	"github.com/flanksource/konfig-manager/pkg"
	"github.com/spf13/cobra"
)

var (
	input, output, outputType, configFile string
)

// GenerateCmd represents the base command when called without any subcommands
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates application-properties based on the hierarchy and config provided",
	RunE: func(cmd *cobra.Command, args []string) error {
		resources, err := pkg.ReadResources(input)
		if err != nil {
			return err
		}
		logger.Infof("%d resources found from %s", len(resources), input)

		for _, r := range resources {
			logger.Debugf("%s", kommons.GetName(r.Item))
		}
		buf, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		var config pkg.Config
		err = yaml.Unmarshal(buf, &config)
		if err != nil {
			return err
		}

		for _, name := range config.Applications {
			logger.Infof("[%s.properties]", name)
			hierarchy, err := pkg.GetHierarchy(configFile, name)
			if err != nil {
				return err
			}

			file := hierarchy.GeneratePropertiesFile(resources)
			if outputType == "stdout" {
				fmt.Println(file)
			} else {
				filePath, err := text.Template(output, map[string]string{"name": name})
				if err != nil {
					return err
				}
				if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
					return err
				}
				if err := ioutil.WriteFile(filePath, []byte(file), 0644); err != nil {
					return err
				}
			}

		}

		return nil
	},
}

func init() {
	GenerateCmd.Flags().StringVarP(&input, "input", "i", "-", "input of yaml dump (default '-' is stdin)")
	GenerateCmd.Flags().StringVarP(&configFile, "config", "c", "config.yml", "path to config file containing hierarchy")
	GenerateCmd.Flags().StringVarP(&outputType, "output-type", "", "stdout", "type of output: can be one of stdout, properties")
	GenerateCmd.Flags().StringVarP(&output, "output-path", "", "properties/{{.name}}.properties", "output path")
}
