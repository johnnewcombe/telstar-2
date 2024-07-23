package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
	"strconv"
)

var publishFrames = &cobra.Command{
	Use:   "purge-frame",
	Short: "Deletes a page, i.e. a frame and all follow-on frames.",
	Long: `
Deletes a page, i.e. a frame and all follow-on frames, including zero page routed pages.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			count    int
			respData network.ResponseData
		)

		// FIXME Add code here
		// FIXME Not implemented at server ??

		result := map[string]string{
			"Saved": strconv.Itoa(count),
		}

		stdOut(cmd, respData, result)
		return nil
	},
}
