package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
)

var deleteUser = &cobra.Command{
	Use:   "delete-user",
	Short: "Deletes a user from the currently logged in system.",
	Long: `
Deletes a user from the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl   string
			userId   string
			token    string
			err      error
			respData network.ResponseData
		)

		//load token - don't want to report errors here as we want an unauthorised status code to be returned
		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if userId, err = cmd.Flags().GetString("user-id"); err != nil {
			return err
		}

		apiUrl += "/user/" + userId

		respData, err = network.Delete(apiUrl, token)
		if err != nil {
			return err
		}

		stdOut(cmd, respData, nil)

		return nil

	},
}
