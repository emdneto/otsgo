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
	"io/ioutil"
	"log"
	"os"

	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var res bool
var auth client.AuthYaml

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Perform basic http auth and store credentials",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		res = client.GetStatus(AuthInfo)
		if res {
			auth = client.AuthYaml{
				Username: username,
				Password: password,
				Enabled:  true,
			}
		} else {
			auth = client.AuthYaml{
				Username: "",
				Password: "",
				Enabled:  false,
			}
		}
		fmt.Printf("WARNING! Your password will be stored unencrypted in %s\n", viper.ConfigFileUsed())
		fmt.Printf("\n")
		fmt.Printf("Login Succeeded\n")
		yamlData, err := yaml.Marshal(&auth)
		if err != nil {
			fmt.Printf("Error while Marshaling. %v", err)
		}

		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		fileName := fmt.Sprintf("%s/.otsgo.yaml", dirname)
		err = ioutil.WriteFile(fileName, yamlData, 0644)
		if err != nil {
			panic("Unable to write data into the file")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringP("username", "u", "", "Username")
	loginCmd.PersistentFlags().StringP("password", "p", "", "Password")
	//loginCmd.PersistentFlags().BoolP("password-stdin", "", false, "Take the API Token from stdin")
}
