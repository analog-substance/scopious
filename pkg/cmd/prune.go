package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// PruneCmd represents the check command
var PruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune excluded scope items from input",
	Long: `Prune excluded scope items from input

cat urls.txt | scopious prune
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")
		scope := scoperInstance.GetScope(scopeName)

		scopePrinted := map[string]bool{}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			scopeLine := scanner.Text()
			if scope.IsInScope(scopeLine) {
				if _, ok := scopePrinted[scopeLine]; !ok {
					scopePrinted[scopeLine] = true
					fmt.Println(scopeLine)
				}
			}
		}

		if scanner.Err() != nil {
			log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
		}
	},
}

func init() {
	RootCmd.AddCommand(PruneCmd)
}
