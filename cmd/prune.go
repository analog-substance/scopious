package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
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
		prunedScope := scope.Prune(scopeToCheck...)

		//sort.Strings(prunedScope)

		for _, scopeLine := range prunedScope {
			fmt.Println(scopeLine)
		}
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pruneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pruneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
