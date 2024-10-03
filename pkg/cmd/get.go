package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// GetCmd represents the add command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get scope things",
	Long: `get scope things
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")

		ipv4, _ := cmd.Flags().GetBool("ipv4")
		ipv6, _ := cmd.Flags().GetBool("ipv6")
		domain, _ := cmd.Flags().GetBool("domain")
		exclude, _ := cmd.Flags().GetBool("exclude")

		if ipv4 {
			fmt.Println(scoperInstance.GetScopeIPv4Path(scopeName))
		}

		if ipv6 {
			fmt.Println(scoperInstance.GetScopeIPv6Path(scopeName))
		}

		if domain {
			fmt.Println(scoperInstance.GetScopeDomainsPath(scopeName))
		}

		if exclude {
			fmt.Println(scoperInstance.GetScopeExcludePath(scopeName))
		}

	},
}

func init() {
	RootCmd.AddCommand(GetCmd)
	GetCmd.Flags().BoolP("ipv4", "4", false, "Get IPv4 file path")
	GetCmd.Flags().BoolP("ipv6", "6", false, "Get IPv6 file path")
	GetCmd.Flags().BoolP("domain", "d", false, "Get domains file path")
	GetCmd.Flags().BoolP("exclude", "x", false, "Get exclude file path")
}
