package cmd

import "github.com/spf13/cobra"

var deleteUser = &cobra.Command{
	Use:   "delete-user",
	Short: "Deletes a user from the currently logged in system.",
	Long: `
Deletes a user from the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
