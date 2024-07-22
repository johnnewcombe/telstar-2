package cmd

import "github.com/spf13/cobra"

var deleteFrame = &cobra.Command{
	Use:   "delete-frame",
	Short: "Deletes a single frame from the currently logged in system.",
	Long: `
Deletes a single frame from the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
