package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var getFrame = &cobra.Command{
	Use:   "get-frame",
	Short: "Returns a single frame from the currently logged in system.",
	Long: `
Returns a single frame from the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl   string
			pageId   string
			primary  bool
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

		if !utils.IsValidPageId(pageId) {
			return errors.New("invalid frame id")
		}

		apiUrl = apiUrl + "/frame/" + pageId
		if primary {
			apiUrl += "?db=primary"
		}

		respData, err = network.Get(apiUrl, token)
		if err != nil {
			return err
		}

		fmt.Println(respData.Body)
		return nil

	},
}
