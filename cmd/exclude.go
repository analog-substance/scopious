package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// excludeCmd represents the block command
var excludeCmd = &cobra.Command{
	Use:   "exclude",
	Short: "Add an item to the exclude list",
	Long: `Add items to to the exclude list.

Sometimes not every subdomain underneath a domain or IP address
in a CIDR is in scope.

	scopious exclude admin.example.com
`,
	Run: func(cmd *cobra.Command, args []string) {
		shouldList, _ := cmd.Flags().GetBool("list")
		scopeName, _ := cmd.Flags().GetString("scope")
		scope := scoperInstance.GetScope(scopeName)

		if shouldList {
			for excluded, _ := range scope.Excludes {
				fmt.Println(excluded)
			}
			return
		}
		if len(args) > 0 {
			scope.AddExclude(args...)
		} else {
			// no args, lets read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				scopeLine := scanner.Text()
				log.Println(scopeLine)
				scope.AddExclude(scopeLine)
			}

			if scanner.Err() != nil {
				log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
			}
		}
		scoperInstance.Save()
	},
}

func init() {
	rootCmd.AddCommand(excludeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// excludeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	excludeCmd.Flags().BoolP("list", "l", false, "List excluded scope")
}
