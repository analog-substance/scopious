package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

// ipsCmd represents the ips command
var ipsCmd = &cobra.Command{
	Use:   "ips",
	Short: "List IP addresses in scope",
	Long: `List IP addresses in scope.

Show in scope ips
	scopious ips

Expand CIDRs and remove excluded ips
	scopious ips -x
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")
		shouldExpand, _ := cmd.Flags().GetBool("expand")
		all, _ := cmd.Flags().GetBool("all")
		scope := scoperInstance.GetScope(scopeName)

		var scopeStrings []string
		if shouldExpand {
			scopeStrings = scope.AllExpanded(all)
			sort.Strings(scopeStrings)
		} else {
			scopeStrings = scope.AllIPs()
		}

		for _, ip := range scopeStrings {
			fmt.Println(ip)
		}
	},
}

func init() {
	rootCmd.AddCommand(ipsCmd)
	ipsCmd.Flags().BoolP("expand", "x", false, "Expand CIDRS and remove excluded things")
	ipsCmd.PersistentFlags().BoolP("all", "a", false, "show all addreses, even network and broadcast")
}
