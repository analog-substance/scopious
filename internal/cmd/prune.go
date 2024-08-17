package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// pruneCmd represents the check command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune excluded scope items from input",
	Long: `Prune excluded scope items from input

cat urls.txt | scopious prune
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")
		all, _ := cmd.Flags().GetBool("all")
		scope := scoperInstance.GetScope(scopeName)

		scopeToCheck := []string{}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			scopeLine := scanner.Text()
			scopeToCheck = append(scopeToCheck, scopeLine)
		}

		if scanner.Err() != nil {
			log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
		}
		prunedScope := scope.Prune(all, scopeToCheck...)

		//sort.Strings(prunedScope)

		for _, scopeLine := range prunedScope {
			fmt.Println(scopeLine)
		}
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
	pruneCmd.PersistentFlags().BoolP("all", "a", false, "show all addreses, even network and broadcast")
}
