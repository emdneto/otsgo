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
	"fmt"

	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
)

// burnCmd represents the burn command
var burnCmd = &cobra.Command{
	Use:   "burn",
	Short: "Burn a secret that has not been read yet",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments. You can only have one which is the METADATA_KEY.")
		} else {
			secret := args[0]
			postBody := client.SecretBody{
				Secret: secret,
			}
			fmt.Println(secret)
			client.BurnSecret(AuthInfo, postBody)
		}
	},
}

func init() {
	rootCmd.AddCommand(burnCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// burnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// burnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
