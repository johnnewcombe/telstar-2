package cmd

import (
	"fmt"
	"github.com/johnnewcombe/telstar-util/globals"
	"github.com/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
)

var login = &cobra.Command{
	Use:   "login",
	Short: "Logs into a system.",
	Long: `
Logs into a system. A successful login stores a token to the local filesystem
this token is used on all subsequent commands that need authentication. The
token expires after 10 minutes of inactivity.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl   string
			userId   string
			password string
			err      error
			respData network.ResponseData
		)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if userId, err = cmd.Flags().GetString("user-id"); err != nil {
			return err
		}

		if password, err = cmd.Flags().GetString("password"); err != nil {
			return err
		}

		// specific case of put that returns a token
		data := fmt.Sprintf("{\"user-id\": \"%s\", \"password\": \"%s\"}", userId, password)

		if respData, err = network.Put(apiUrl+"/login", data, ""); err != nil {
			return err
		}

		if err = saveText(globals.TOKENFILE, respData.Token); err != nil {
			return (err)
		}

		//fmt.Printf(globals.Response, respData.Status)
		stdOut(cmd, respData, nil)

		return nil
	},
}
