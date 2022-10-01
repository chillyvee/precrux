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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chillyvee/precrux/snitch"
	"github.com/spf13/cobra"
)

var ()

func init() {
	snitchCmd.AddCommand(snitchStartCmd)
	snitchCmd.AddCommand(snitchAddChaserCmd)
	rootCmd.AddCommand(snitchCmd)
}

// versionCmd represents the version command
var snitchCmd = &cobra.Command{
	Use:   "snitch",
	Short: "Snitch mode for precrux",
}

var snitchAddChaserCmd = &cobra.Command{
	Use:   "add [chaser-name]",
	Short: "Add remote chaser by name",
	Long: "Add remote chaser by name\n\n" +
		"[chaser-name] is the name of the chaser server.  Can be a hostname or any other word like 'west-1' or 'server-2'",
	Args:         cobra.RangeArgs(1, 1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Paste certificate printerd on chaser.  Input a blank line to finish.")
		fmt.Println("First Line Should Start with -----BEGIN CERTIFICATE-----")
		fmt.Println("")

		var textBuffer strings.Builder
		textBuffer.Grow(10240)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if len(scanner.Text()) == 0 {
				break
			}
			textBuffer.WriteString(scanner.Text() + "\n")
		}
		certstring := textBuffer.String()

		certfile := fmt.Sprintf("chaser_%s.crt", args[0])
		if err := os.WriteFile(certfile, []byte(certstring), 600); err != nil {
			panic(err)
		}

		fmt.Printf("Chaser %s certificate saved to %s\n", args[0], certfile)
		return nil
	},
}

var snitchStartCmd = &cobra.Command{
	Use:   "start [snitch-name] [ip:port]",
	Short: "Start Snitch mode for precrux",
	Long: "Start Chaser mode for precrux\n\n" +
		"[snitch-name] is the name of the horcrux shard server.  Can be a hostname or any other word like 'west-1' or 'server-2'",
	Args:         cobra.RangeArgs(2, 2),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &snitch.Snitch{ChaserName: args[0], ChaserAddress: args[1]}
		c.Start()
		return nil
	},
}
