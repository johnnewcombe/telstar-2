package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	k_apiUrl      = "URL of the API to be accessed."
	k_userId      = "User ID."
	k_userIdAdd   = "User ID to add or update."
	k_password    = "Password."
	k_name        = "User's display name."
	k_pageId      = "Page ID."
	k_primary     = "Specify the primary database."
	k_source      = "Directory or GIT repository containing the json files to upload."
	k_destination = "Directory where the json files should be stored."
	k_fileName    = "Filename of json file."
	k_json        = "Output results in json"
	k_admin       = "Sets the admin status of the new user, default on non-admin."
	k_api         = "Sets the api access status of the new user, default on no-access."
	k_editor      = "Sets the editor status of the new user, default on non-editor."
	k_basepage    = "Sets the base page. The base page determines where in the tree that frames can be added. Default is 999999999"
	k_purge       = "Deletes all of the associated frames including 'zero page' extension frames"
)

func init() {

	rootCmd.AddCommand(addFrame)
	rootCmd.AddCommand(addFrames)
	rootCmd.AddCommand(addUser)
	rootCmd.AddCommand(deleteFrame)
	rootCmd.AddCommand(deleteUser)
	rootCmd.AddCommand(getFrame)
	rootCmd.AddCommand(getFrames)
	rootCmd.AddCommand(getStatus)
	rootCmd.AddCommand(login)
	rootCmd.AddCommand(purgeFrames)

	rootCmd.PersistentFlags().BoolP("json", "j", false, k_json)
	rootCmd.PersistentFlags().String("url", "", k_apiUrl)

	// Add Frame
	addFrame.PersistentFlags().StringP("filename", "f", "", k_fileName)
	addFrame.PersistentFlags().Bool("primary", false, k_primary)

	// Add Frames
	addFrames.PersistentFlags().StringP("source", "s", "", k_source)
	addFrames.PersistentFlags().Bool("primary", false, k_primary)

	// Add User
	//addUser.PersistentFlags().String("url", "", k_apiUrl)
	addUser.PersistentFlags().StringP("user-id", "u", "", k_userIdAdd)
	addUser.PersistentFlags().StringP("password", "p", "", k_password)
	addUser.PersistentFlags().StringP("name", "n", "", k_name)
	addUser.PersistentFlags().BoolP("admin", "a", false, k_admin)
	addUser.PersistentFlags().Bool("api-access", false, k_api)
	addUser.PersistentFlags().Bool("editor", false, k_editor)
	addUser.PersistentFlags().IntP("base-page", "b", 999999999, k_basepage)

	// Delete Frame
	deleteFrame.PersistentFlags().String("page-id", "", k_pageId)
	deleteFrame.PersistentFlags().Bool("primary", false, k_primary)
	deleteFrame.PersistentFlags().Bool("purge", false, k_purge)

	// Delete User
	deleteUser.PersistentFlags().StringP("user-id", "u", "", k_userId)

	// Get Frame <url> <page id> [primary|secondary]"
	getFrame.PersistentFlags().String("page-id", "", k_pageId)
	getFrame.PersistentFlags().Bool("primary", false, k_primary)

	// Get Frames
	getFrames.PersistentFlags().StringP("destination", "d", "", k_destination)
	getFrames.PersistentFlags().Bool("primary", false, k_primary)

	// Get Status
	// none

	// Login
	login.PersistentFlags().StringP("user-id", "u", "", k_userId)
	login.PersistentFlags().StringP("password", "p", "", k_password)

	// Purge Frames
	purgeFrames.PersistentFlags().String("page-id", "", k_pageId)
	purgeFrames.PersistentFlags().Bool("primary", false, k_primary)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error()+".")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "telstar-util",
	Short: "Utility program for interacting with the Telstar Server API. Version: " + globals.Version,
	Long:  `Utility program for interacting with the Telstar Server API.`,
}
