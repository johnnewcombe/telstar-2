package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"fmt"
	"github.com/spf13/cobra"
)

var addUser = &cobra.Command{
	Use:   "add-user",
	Short: "Adds/updates a user in the currently logged in system.",
	Long: `
Adds or update a user in the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// template to create the json data, easier than marshalling a type
		const dataTemplate = "{\"user-id\":\"%s\",\"password\":\"%s\",\"admin\":%t,\"api-access\":%t,\"editor\":%t,\"name\":\"%s\",\"base-page\":%d}"

		var (
			respData network.ResponseData
			apiUrl   string
			err      error
			userId   string
			password string
			name     string
			token    string

			// TODO these could be set by cmd line switches
			basePage  int
			admin     bool
			editor    bool
			apiAccess bool
		)

		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if userId, err = cmd.Flags().GetString("user-id"); err != nil {
			return err
		}

		if password, err = cmd.Flags().GetString("password"); err != nil {
			return err
		}

		if name, err = cmd.Flags().GetString("name"); err != nil {
			return err
		}

		if basePage, err = cmd.Flags().GetInt("base-page"); err != nil {
			return err
		}

		if admin, err = cmd.Flags().GetBool("admin"); err != nil {
			return err
		}

		if apiAccess, err = cmd.Flags().GetBool("api-access"); err != nil {
			return err
		}

		if editor, err = cmd.Flags().GetBool("editor"); err != nil {
			return err
		}

		apiUrl = apiUrl + "/user"

		data := fmt.Sprintf(dataTemplate, userId, password, admin, apiAccess, editor, name, basePage)

		respData, err = network.Put(apiUrl, data, token)
		if err != nil {
			return err
		}

		// TODO the error should come from the api
		//if !isValidUserId(userId) {
		//	return errors.New("bad userid")
		//}

		stdOut(cmd, respData, nil)

		return nil
	},
}