package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registrantCmd represents the registrant command
var registrantCmd = &cobra.Command{
	Use:   "registrant",
	Short: "Add or remove registrants",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Example: `klinkregistry registrant add --email admin@example.com --admin=true`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		switch args[0] {
		case "add":
			fmt.Println("Adding Registrant")
		case "del":
			fmt.Println("Deleting Registrant")
		default:
			fmt.Printf("Unknown command: %s\n", args[0])
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(registrantCmd)
}
