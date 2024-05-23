package cmd

import (
	"bufio"
	"fmt"
	"github.com/analog-substance/scopious/pkg/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var expandCmd = &cobra.Command{
	Use:   "expand",
	Short: "Expand CIDRs",
	Long: `Expand CIDRs. For example:

	cat customer-supplied.txt | scopious expand

	scopious expand 10.0.0.0/22
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			for _, scopeLine := range args {
				processScopeLine(scopeLine)
			}
		} else {
			// no args, lets read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				scopeLine := scanner.Text()
				processScopeLine(scopeLine)
			}

			if scanner.Err() != nil {
				log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
			}
		}
	},
}

func processScopeLine(scopeLine string) {
	if strings.Contains(scopeLine, "/") {
		// perhaps we have a CIDR
		ips, err := utils.GetAllIPs(scopeLine)
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
}
