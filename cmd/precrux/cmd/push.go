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
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push [chain-name] [chaser-name]",
	Short: "Push [chain-name] information to [chaser-name]",
	Long: "Push chain-name information to chaser-name\n\n" +
		"[chain-name] is the name of the chain.  Example 'juno'\n" +
		"[chaser-name] is the name of the horcrux shard server.  Can be a hostname or any other word like 'west-1' or 'server-2'",
	Args:         cobra.RangeArgs(2, 2),
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := &snitch.Snitch{ChainName: args[0], ChaserName: args[1]}
		s.SendChainToChaser()
		return nil
	},
}
