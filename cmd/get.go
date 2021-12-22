package cmd

import (
	"fmt"

	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret, metadata or recent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify a resource like secret, meta or recent")
	},
}

var GetSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Retrieve a Secret",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments. You can only have one which is the secret URL or SECRET_KEY")
		} else {
			secret := args[0]
			pp, _ := cmd.Flags().GetString("passphrase")
			postBody := client.SecretBody{
				Secret:     secret,
				Passphrase: pp,
			}
			client.GetSecret(AuthInfo, postBody)
		}

	},
}

var GetMetadataCmd = &cobra.Command{
	Use:   "meta",
	Short: "Retrieve secret associated metadata",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments. You can only have one which is the METADATA_KEY.")
		} else {
			secret := args[0]
			postBody := client.SecretBody{
				Secret: secret,
			}
			fmt.Println(secret)
			client.GetMetadata(AuthInfo, postBody)
		}
	},
}

var GetRecentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Retreive a list of recent metadata.",
	Run: func(cmd *cobra.Command, args []string) {
		client.GetRecent(AuthInfo, client.SecretBody{})

	},
}

func init() {
	getCmd.AddCommand(GetSecretCmd)
	getCmd.AddCommand(GetMetadataCmd)
	getCmd.AddCommand(GetRecentCmd)
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
