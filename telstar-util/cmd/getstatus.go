package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"fmt"
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
			url      string
			respData network.ResponseData
			err      error
		)

		// get the url to be checked
		if url, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if respData, err = network.Get(url, ""); err != nil {
			return err
		}

		fmt.Printf(globals.Response, respData.Status)
		return nil
	},
}
