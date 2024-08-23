package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// DomainsCmd represents the domains command
var DomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Print out in scope domains",
	Long: `Print out in scope domains. For example:

Print all domains in scope:
	scopious domains

Print in scope root domains:
	scopious domains -r
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")
		showRootDomains, _ := cmd.Flags().GetBool("root-domains")
		scope := scoperInstance.GetScope(scopeName)

		var domains []string

		if showRootDomains {
			domains = scope.RootDomains()
		} else {
			domains = scope.AllDomains()
		}

		for _, domain := range domains {
			fmt.Println(domain)
		}
	},
}

func init() {
	RootCmd.AddCommand(DomainsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	DomainsCmd.Flags().BoolP("root-domains", "r", false, "Show only root domains")
}
