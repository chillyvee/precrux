/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/chillyvee/precrux/snitch"
	"github.com/spf13/cobra"
)

var ()

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate [chain-name]",
	Short: "Generate files to send to remote Chasers",
	Long: "Generate files to send to remote Chasers\n\n" +
		"Expects [chain-name]/priv_validator_key.json\n" +
		"Expects [chain-name]/precrux.yaml to configure chain",
	Args:         cobra.RangeArgs(1, 1),
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := &snitch.Snitch{ChainName: args[0]}
		s.Generate()
		s.ImportAndWriteSignState()
		return nil
	},
}
