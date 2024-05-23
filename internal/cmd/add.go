package cmd

import (
	"bufio"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add items to scope",
	Long: `Add items to scope. For example:

	cat customer-supplied.txt | scopious add

	scopious add -i internal 10.0.0.0/22
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")
		scope := scoperInstance.GetScope(scopeName)

		if len(args) > 0 {
			scope.Add(args...)
		} else {
			// no args, lets read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				scopeLine := scanner.Text()
				scope.Add(scopeLine)
			}

			if scanner.Err() != nil {
				log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
			}
		}

		scoperInstance.Save()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//addCmd.PersistentFlags().StringP("scope", "s", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
