package cmd

import (
	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Current status of the system",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client.GetStatus(AuthInfo)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
