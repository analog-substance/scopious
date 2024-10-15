package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
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
		totals, _ := cmd.Flags().GetBool("totals")
		withSuffix, _ := cmd.Flags().GetString("suffix")
		scope := scoperInstance.GetScope(scopeName)

		var domains []string

		if showRootDomains {
			domains = scope.RootDomains()
			if totals {
				totalMap := map[string]int{}
				allDomains := scope.AllDomains()
				for _, rootDomain := range domains {
					for _, domain := range allDomains {
						if strings.HasSuffix(domain, rootDomain) {
							totalMap[rootDomain]++
						}
					}
				}
				for rootDomain, count := range totalMap {
					fmt.Println(rootDomain, count)
				}
			}
		} else {
			domains = scope.AllDomains()
		}

		for _, domain := range domains {
			//if totals {
			//
			//} else {
			if withSuffix != "" {
				if strings.HasSuffix(domain, withSuffix) {
					fmt.Println(domain)
				}
			} else {
				fmt.Println(domain)
			}
			//}
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
	DomainsCmd.Flags().StringP("suffix", "S", "", "Show only domains with suffix")
	DomainsCmd.Flags().BoolP("totals", "t", false, "Show totals for root domains and suffix")
}
