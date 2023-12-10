package cmd

import (
	"fmt"
	"github.com/analog-substance/scopious/pkg/scopious"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var scoperInstance *scopious.Scoper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scopious",
	Short: "Manage scope for your network based projects",
	Long: `Scoper can help you manage the scope of network projects by:

- Automatically detecting and separating IP addresses or domains
- Ensuring an item is in the scope of your engagement
- Keep track of multiple scope for your engagement

To use, simply supply your scope as arguments to scopious add

	scopious add example.com example.net 203.0.113.0/24
	cat scope.txt | scopious add

By default scope is stored in ./scope/external/. This scan be changed by specifying -s

	scopious add -s internal evil.corp internal.corpdev 10.0.0.1/24

You can exclude things from scope as well

	scope excluded admin.example.com 203.0.113.0/29

Scoper can validate items are in scope

	cat maybe-inscope.txt | scopious prune > inscope.txt

Need to view your scope data, scopious can show you all your scope in various ways

List in scope domains
	scopious domains

list in scope root domains
	scopious domains -r

list in scope ips
	scopious ips

expand cidrs and remove excluded things
	scopious ips -x

list excluded things
	scopious exclude -l




`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		scopeDir := viper.GetString("scope-dir")
		scoperInstance = scopious.FromPath(scopeDir)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scopious.yaml)")

	rootCmd.PersistentFlags().String("scope-dir", "scope", "where scope files are located.")
	rootCmd.PersistentFlags().StringP("scope", "s", scopious.DefaultScope, "Scope name")

	//rootCmd.PersistentFlags().String("domains-file", "scope-domains.txt", "where in-scope domains are located.")
	//rootCmd.PersistentFlags().String("ips-file", "scope-ips.txt", "where in-scope IP addresses are located.")
	//rootCmd.PersistentFlags().String("ignore-domains", "ignore-scope-domains.txt", "where out-of-scope IP addresses are located.")
	//rootCmd.PersistentFlags().String("ignore-ips", "ignore-scope-ips.txt", "where out-of-scope domains addresses are located.")

	viper.BindPFlag("scope-dir", rootCmd.PersistentFlags().Lookup("scope-dir"))
	//viper.BindPFlag("ips-file", rootCmd.PersistentFlags().Lookup("ips-file"))
	//viper.BindPFlag("ignore-domains", rootCmd.PersistentFlags().Lookup("ignore-domains"))
	//viper.BindPFlag("ignore-ips", rootCmd.PersistentFlags().Lookup("ignore-ips"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".scopious" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".scopious")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
