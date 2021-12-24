package cmd

import (
	"bufio"
	"os"

	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share or generate a random secret",
	Run: func(cmd *cobra.Command, args []string) {
		frmStdin, _ := cmd.Flags().GetBool("from-stdin")
		gen, _ := cmd.Flags().GetBool("generate")
		secret, _ := cmd.Flags().GetString("secret")
		if frmStdin {
			var multiline string
			var reader = bufio.NewReader(os.Stdin)
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				multiline += scanner.Text() + "\n"
			}
			secret = multiline
		}

		pp, _ := cmd.Flags().GetString("passphrase")
		reci, _ := cmd.Flags().GetString("recipient")
		ttl, _ := cmd.Flags().GetInt("ttl")
		postBody := client.SecretBody{
			Secret:     secret,
			Passphrase: pp,
			Recipient:  reci,
			Ttl:        ttl,
		}
		client.CreateSecret(AuthInfo, postBody, gen)

	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.PersistentFlags().StringP("secret", "s", "", "the secret value which is encrypted before being stored. There is a maximum length based on your plan that is enforced (1k-10k)")
	shareCmd.PersistentFlags().StringP("passphrase", "p", "", "a string that the recipient must know to view the secret. This value is also used to encrypt the secret and is bcrypted before being stored so we only have this value in transit.")
	shareCmd.PersistentFlags().StringP("recipient", "r", "", "an email address. We will send a friendly email containing the secret link (NOT the secret itself).")
	shareCmd.PersistentFlags().IntP("ttl", "t", 604800, "the maximum amount of time, in seconds, that the secret should survive (i.e. time-to-live). Once this time expires, the secret will be deleted and not recoverable.")
	shareCmd.Flags().BoolP("from-stdin", "f", false, "Read from stdin")
	shareCmd.Flags().BoolP("generate", "g", false, "Generate a short, unique secret. This is useful for temporary passwords, one-time pads, salts, etc.")
}
