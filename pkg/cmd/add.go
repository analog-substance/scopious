package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add items to scope",
	Long: `Add items to scope unless it has been excluded via scopious exclude. For example:

	cat customer-supplied.txt | scopious add

	scopious add -i internal 10.0.0.0/22
`,
	Run: func(cmd *cobra.Command, args []string) {
		scopeName, _ := cmd.Flags().GetString("scope")
		all, _ := cmd.Flags().GetBool("all")
		scope := scoperInstance.GetScope(scopeName)

		if len(args) > 0 {
			scope.Add(all, args...)
		} else {
			// no args, lets read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				scopeLine := scanner.Text()
				scope.Add(all, scopeLine)
			}

			if scanner.Err() != nil {
				log.Printf("STDIN scanner encountered an error: %s", scanner.Err())
			}
		}

		scoperInstance.Save()
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
	AddCmd.PersistentFlags().BoolP("all", "a", false, "show all addresses, even network and broadcast")
}
