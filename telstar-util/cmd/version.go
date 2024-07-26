package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"fmt"
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

		fmt.Printf("\ntelstar-util %s\n\n", globals.Version)

		return nil
	},
}
