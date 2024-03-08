package cmd

import (
	"fmt"
	"os"

	"github.com/emdneto/ots/client"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var AuthInfo client.Auth
var (
	// VERSION is set during build
	VERSION string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ots",
	Short: "A simple CLI and API client for One-Time Secret",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Println(cmd.Use + " " + VERSION)
		} else {
			cmd.Help()
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	VERSION = v
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.otsgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Displays current version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".otsgo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	var AuthEnabled bool = false // default to false
	var AuthUsername string
	var AuthPassword string

	AuthUsername = viper.GetString("OTS_USER")
	AuthPassword = viper.GetString("OTS_TOKEN")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		AuthUsername = viper.GetString("otsUser")
		AuthPassword = viper.GetString("otsToken")

	}

	if len(AuthUsername) != 0 && len(AuthPassword) != 0 {
		AuthEnabled = true
		AuthInfo = client.Auth{
			Username: AuthUsername,
			Password: AuthPassword,
			Enabled:  AuthEnabled,
		}
	}
}
