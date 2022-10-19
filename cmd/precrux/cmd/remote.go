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
	remoteCmd.AddCommand(remoteAddCmd)
	rootCmd.AddCommand(remoteCmd)
}

// versionCmd represents the version command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote chaser admin",
}

var remoteAddCmd = &cobra.Command{
	Use:   "add [chaser-name] ip:port",
	Short: "Add remote chaser by name",
	Long: "Add remote chaser by name\n\n" +
		"[chaser-name] is the name of the chaser server.  Can be a word like 'west-1' or 'server-2'" +
		"ip:port is the ip and port of the chaser",
	Args:         cobra.RangeArgs(2, 2),
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := &snitch.Snitch{}
		//s.ReadChaserCertificate(args[0])
		s.AddChaserProfile(args[0], args[1])
		return nil
	},
}
