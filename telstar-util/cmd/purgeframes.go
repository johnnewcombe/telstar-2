package cmd

import "github.com/spf13/cobra"

var purgeFrames = &cobra.Command{
	Use:   "purge-frame",
	Short: "Deletes a page, i.e. a frame and all follow-on frames.",
	Long: `
Deletes a page, i.e. a frame and all follow-on frames, including zero page routed pages.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
