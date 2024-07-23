package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
)

var deleteFrame = &cobra.Command{
	Use:   "delete-frame",
	Short: "Deletes a single frame from the currently logged in system.",
	Long: `
Deletes a single frame from the currently logged in system. See the login command.`,
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

		if pageId, err = cmd.Flags().GetString("page-id"); err != nil {
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

		respData, err = network.Delete(apiUrl, token)
		if err != nil {
			return err
		}

		stdOut(cmd, respData, nil)

		return nil
	},
}
