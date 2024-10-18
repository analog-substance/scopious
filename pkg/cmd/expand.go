package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/analog-substance/scopious/pkg/utils"
	"github.com/spf13/cobra"
)

var ExpandCmd = &cobra.Command{
	Use:   "expand",
	Short: "Expand CIDRs",
	Long: `Expand CIDRs. For example:

	cat customer-supplied.txt | scopious expand

	scopious expand 10.0.0.0/22
`,
	Run: func(cmd *cobra.Command, args []string) {

		all, _ := cmd.Flags().GetBool("all")
		public, _ := cmd.Flags().GetBool("public")
		private, _ := cmd.Flags().GetBool("private")

		if public && private {
			// silly, that is the same as the default...
			public = false
			private = false
		}

		if len(args) > 0 {
			for _, scopeLine := range args {
				processScopeLine(scopeLine, all, public, private)
			}
		} else {
			// no args, lets read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				scopeLine := scanner.Text()
				processScopeLine(scopeLine, all, public, private)
			}

			if scanner.Err() != nil {
				log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
			}
		}
	},
}

func processScopeLine(scopeLine string, all, public, private bool) {
	if strings.Contains(scopeLine, "/") {
		// perhaps we have a CIDR
		ips, err := utils.GetAllIPs(scopeLine, all)
		if err != nil {
			//log.Println("error processing cidr", err)
			return
		}

		for _, ip := range ips {
			if (public && !ip.IsPrivate()) || (private && ip.IsPrivate()) || (!public && !private) {
				fmt.Println(ip.String())
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(ExpandCmd)
	ExpandCmd.PersistentFlags().BoolP("all", "a", false, "show all addresses, even network and broadcast")
	ExpandCmd.PersistentFlags().Bool("public", false, "Only return public IPs")
	ExpandCmd.PersistentFlags().Bool("private", false, "Only return private IPs")
}
