package cmd

import (
	"github.com/emdneto/ots/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret, metadata or recent",
	Args:  cobra.MinimumNArgs(1),
}

var GetSecretCmd = &cobra.Command{
	Use:   "secret [secret]",
	Short: "Retrieve a Secret value",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secret := args[0]
		pp, _ := cmd.Flags().GetString("passphrase")
		postBody := client.SecretBody{
			Secret:     secret,
			Passphrase: pp,
		}
		client.GetSecret(AuthInfo, postBody)
	},
}
var GetMetadataCmd = &cobra.Command{
	Use:   "meta [key]",
	Short: "Retrieve secret associated metadata",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secret := args[0]
		postBody := client.SecretBody{Secret: secret}
		client.GetMetadata(AuthInfo, postBody)
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
}
