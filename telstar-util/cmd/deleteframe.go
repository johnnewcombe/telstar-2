package cmd

import (
	"github.com/johnnewcombe/telstar-util/globals"
	"github.com/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
)

var deleteFrame = &cobra.Command{
	Use:   "delete-frame",
	Short: "Deletes a single frame from the currently logged in system.",
	Long: `
Deletes a single frame from the currently logged in system. See the login
command.

This command can perform purging of data, for example, if frames 101a, 101b,
101c and 101d existed in the system and the command to deleted frame 101b
was executed, setting perge to true, frmaes 101c and 101d would also be
removed from the system. In other words, all ‘follow on’ frames would be
removed from the system. this extends to zero page routed frames, see below.

In cases where a page needs more than 26 frames such as a large Telesoftware
program or where a large number of news articles extends beyond frame z, a
process of Zero Page Routing takes place. For example if a news article
starting on frame 222z needed a continuation frame, the frame 2220a would be
used. If articles continued to frame 2220z and needed further continuation 
frames, frame 22200a would be used and so on to the maximum page number
length.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl   string
			pageId   string
			primary  bool
			purge    bool
			token    string
			err      error
			respData network.ResponseData
		)

		//load token - don't want to report errors here as we want an unauthorised status code to be returned
		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}
		if pageId, err = cmd.Flags().GetString("frame-id"); err != nil {
			return err
		}
		if primary, err = cmd.Flags().GetBool("primary"); err != nil {
			return err
		}
		if purge, err = cmd.Flags().GetBool("purge"); err != nil {
			return err
		}

		apiUrl += "/frame/" + pageId

		if primary {
			apiUrl += "?db=primary"
			if purge {
				apiUrl += "&purge=true"
			}
		} else {
			if purge {
				apiUrl += "?purge=true"
			}
		}

		if respData, err = deleteSingleFrame(apiUrl, token); err != nil {
			return err
		}

		stdOut(cmd, respData, nil)

		return nil
	},
}

func deleteSingleFrame(apiUrl string, token string) (network.ResponseData, error) {
	var (
		err      error
		respData network.ResponseData
	)

	respData, err = network.Delete(apiUrl, token)
	if err != nil {
		return respData, err
	}

	return respData, nil
}
