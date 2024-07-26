package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Returns the version of the system.",
	Long: `
Returns the version of the system.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			respData network.ResponseData
		)

		respData.SetOK()

		result := map[string]string{
			"Version": globals.Version,
		}
		stdOut(cmd, respData, result)

		return nil
	},
}
