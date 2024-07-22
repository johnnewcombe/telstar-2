package cmd

import "github.com/spf13/cobra"

var addUser = &cobra.Command{
	Use:   "add-user",
	Short: "Adds a user to the currently logged in system.",
	Long: `
Adds a user to the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
