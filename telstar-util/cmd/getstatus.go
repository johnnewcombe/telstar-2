package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
)

var getStatus = &cobra.Command{
	Use:   "get-status",
	Short: "Returns the status of the specified system.",
	Long: `
Returns the status of the specified system. This method does not require a login.

The HTTP Response Code 200 (OK) is returned if the system is OK
`,

	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl   string
			respData network.ResponseData
			err      error
		)

		// get the url to be checked
		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		apiUrl += "/status"

		if respData, err = network.Get(apiUrl, ""); err != nil {
			return err
		}
		result := map[string]string{
			"Client Version": globals.Version,
			"API Version":    respData.Body,
		}
		stdOut(cmd, respData, result)
		return nil
	},
}
