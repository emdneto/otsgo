package cmd

import (
	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
)

// burnCmd represents the burn command
var burnCmd = &cobra.Command{
	Use:   "burn [metadata_key]",
	Short: "Burn a secret that has not been read yet",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secret := args[0]
		postBody := client.SecretBody{
			Secret: secret,
		}
		client.BurnSecret(AuthInfo, postBody)
	},
}

func init() {
	rootCmd.AddCommand(burnCmd)
}
