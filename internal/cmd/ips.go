package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sort"
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
		scope := scoperInstance.GetScope(scopeName)

		var scopeStrings []string
		if shouldExpand {
			scopeStrings = scope.AllExpanded()
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	ipsCmd.Flags().BoolP("expand", "x", false, "Expand CIDRS and remove excluded things")
}
