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
		if len(args) < 1 {
			fmt.Println("Error: must also specify the private METADATA_KEY to burn")
		}	
		secret := args[0]
		postBody := client.SecretBody{
			Secret: secret,
		}
		client.BurnSecret(AuthInfo, postBody)
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
