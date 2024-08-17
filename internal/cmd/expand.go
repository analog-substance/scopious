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

var expandCmd = &cobra.Command{
	Use:   "expand",
	Short: "Expand CIDRs",
	Long: `Expand CIDRs. For example:

	cat customer-supplied.txt | scopious expand

	scopious expand 10.0.0.0/22
`,
	Run: func(cmd *cobra.Command, args []string) {

		all, _ := cmd.Flags().GetBool("all")

		if len(args) > 0 {
			for _, scopeLine := range args {
				processScopeLine(scopeLine, all)
			}
		} else {
			// no args, lets read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				scopeLine := scanner.Text()
				processScopeLine(scopeLine, all)
			}

			if scanner.Err() != nil {
				log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
			}
		}
	},
}

func processScopeLine(scopeLine string, all bool) {
	if strings.Contains(scopeLine, "/") {
		// perhaps we have a CIDR
		ips, err := utils.GetAllIPs(scopeLine, all)
		if err != nil {
			//log.Println("error processing cidr", err)
			return
		}

		for _, ip := range ips {
			fmt.Println(ip.String())
		}
	}
}

func init() {
	rootCmd.AddCommand(expandCmd)
	expandCmd.PersistentFlags().BoolP("all", "a", false, "show all addreses, even network and broadcast")
}
