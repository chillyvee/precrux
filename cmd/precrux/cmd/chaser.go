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
	"github.com/chillyvee/precrux/chaser"
	"github.com/spf13/cobra"
)

var ()

func init() {
	chaserCmd.AddCommand(startCmd)
	rootCmd.AddCommand(chaserCmd)
}

// versionCmd represents the version command
var chaserCmd = &cobra.Command{
	Use:   "chaser",
	Short: "Chaser mode for precrux",
}

var startCmd = &cobra.Command{
	Use:   "start [chaser-name]",
	Short: "Start Chaser mode for precrux",
	Long: "Start Chaser mode for precrux\n\n" +
		"[chaser-name] is the name of the horcrux shard server.  Can be a hostname or any other word like 'west-1' or 'server-2'",
	Args:         cobra.RangeArgs(1, 1),
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &chaser.Chaser{Name: args[0], Port: *port}
		c.Start()
		return nil
	},
}
