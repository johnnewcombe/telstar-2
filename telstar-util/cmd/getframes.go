package cmd

import "github.com/spf13/cobra"

var getFrames = &cobra.Command{
	Use:   "get-frames",
	Short: "Returns multiple frames from the currently logged in system.",
	Long: `
Returns multiple frame from the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
