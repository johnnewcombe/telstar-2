package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"errors"
	"github.com/spf13/cobra"
)

var publishFrame = &cobra.Command{
	Use:   "publish-frame",
	Short: "Publishes frames from the primary database to the secondary.",
	Long: `
The Telstar database consists of two distinct parts, Primary and Secondary.
Typically frames would be updated on the Secondary part and then checked by
viewing them via a Telstar instance configured to point to the Secondary
databse. These could then be published to the Primary once everything was
deemed to be OK.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			respData network.ResponseData
			err      error
			pageId   string
			apiUrl   string
			token    string
		)

		//load token - don't want to report errors here as we want an unauthorised status code to be returned
		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if pageId, err = cmd.Flags().GetString("page-id"); err != nil {
			return err
		}

		apiUrl += "/publish/" + pageId

		if !utils.IsValidPageId(pageId) {
			return errors.New("invalid frame id")
		}

		// FIXME Not implemented at server ??

		respData, err = network.Get(apiUrl, token)
		if err != nil {
			return err
		}

		stdOut(cmd, respData, nil)
		return nil
	},
}
