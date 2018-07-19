package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Add or remove applications",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: application [add,del]")
			return
		}

		switch args[0] {
		case "add":
			fmt.Println("Adding Application")
		case "del":
			fmt.Println("Deleting Application")
		default:
			fmt.Printf("Unknown command: %s\n", args[0])
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(applicationCmd)
}
