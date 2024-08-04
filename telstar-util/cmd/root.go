package cmd

import (
	_ "embed"
	"fmt"
	"github.com/johnnewcombe/telstar-util/globals"
	"github.com/spf13/cobra"
	"os"
)

const (
	k_apiUrl        = "URL of the API to be accessed."
	k_userId        = "User ID."
	k_userIdAdd     = "User ID to add or update."
	k_password      = "Password."
	k_name          = "User's display name."
	k_pageNo        = "Restricts the update to a specific page."
	k_frameId       = "Frame ID."
	k_primary       = "Specify the primary database."
	k_source        = "Directory or GIT repository containing the json files to upload."
	k_destination   = "Directory where the json files should be stored."
	k_fileName      = "Source filename of json file."
	k_json          = "Output results in json"
	k_admin         = "Sets the admin status of the new user, default on non-admin."
	k_api           = "Sets the api access status of the new user, default on no-access."
	k_editor        = "Sets the editor status of the new user, default on non-editor."
	k_basepage      = "Sets the base page. The base page determines where in the tree that frames can be added. Default is 999999999"
	k_purge         = "Deletes all of the associated frames including 'zero page' extension frames"
	k_includeUnsafe = "Include frames such as response frames and gateway frames."
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
	rootCmd.AddCommand(publishFrame)
	rootCmd.AddCommand(addPage)
	rootCmd.AddCommand(version)

	rootCmd.PersistentFlags().String("url", "", k_apiUrl)

	// Get Status
	getStatus.PersistentFlags().BoolP("json", "j", false, k_json)

	// Login
	login.PersistentFlags().StringP("user-id", "u", "", k_userId)
	login.PersistentFlags().StringP("password", "p", "", k_password)
	login.PersistentFlags().BoolP("json", "j", false, k_json)

	// Get Frame <url> <page id> [primary|secondary]"
	getFrame.PersistentFlags().StringP("frame-id", "f", "", k_frameId)
	getFrame.PersistentFlags().Bool("primary", false, k_primary)

	// Get Frames
	getFrames.PersistentFlags().StringP("destination", "d", "", k_destination)
	getFrames.PersistentFlags().Bool("primary", false, k_primary)
	getFrames.PersistentFlags().BoolP("json", "j", false, k_json)

	// Add Frame
	addFrame.PersistentFlags().StringP("source", "s", "", k_fileName)
	addFrame.PersistentFlags().Bool("primary", false, k_primary)
	addFrame.PersistentFlags().BoolP("json", "j", false, k_json)
	addFrame.PersistentFlags().Bool("include-unsafe", false, k_includeUnsafe)

	// Add Frames
	addFrames.PersistentFlags().StringP("source", "s", "", k_source)
	addFrames.PersistentFlags().Bool("primary", false, k_primary)
	addFrames.PersistentFlags().BoolP("json", "j", false, k_json)
	addFrames.PersistentFlags().Bool("include-unsafe", false, k_includeUnsafe)

	// Add Page
	addPage.PersistentFlags().StringP("source", "s", "", k_source)
	addPage.PersistentFlags().Int("page-no", -1, k_pageNo)
	addPage.PersistentFlags().Bool("primary", false, k_primary)
	addPage.PersistentFlags().BoolP("json", "j", false, k_json)
	addPage.PersistentFlags().Bool("include-unsafe", false, k_includeUnsafe)

	// Delete Frame
	deleteFrame.PersistentFlags().String("frame-id", "f", k_frameId)
	deleteFrame.PersistentFlags().Bool("primary", false, k_primary)
	deleteFrame.PersistentFlags().Bool("purge", false, k_purge)
	deleteFrame.PersistentFlags().BoolP("json", "j", false, k_json)

	// Publish Frame
	publishFrame.PersistentFlags().String("frame-id", "f", k_frameId)
	publishFrame.PersistentFlags().BoolP("json", "j", false, k_json)

	// Add User
	addUser.PersistentFlags().StringP("user-id", "u", "", k_userIdAdd)
	addUser.PersistentFlags().StringP("password", "p", "", k_password)
	addUser.PersistentFlags().StringP("name", "n", "", k_name)
	addUser.PersistentFlags().BoolP("admin", "a", false, k_admin)
	addUser.PersistentFlags().Bool("api-access", false, k_api)
	addUser.PersistentFlags().Bool("editor", false, k_editor)
	addUser.PersistentFlags().IntP("base-page", "b", 999999999, k_basepage)
	addUser.PersistentFlags().BoolP("json", "j", false, k_json)

	// Delete User
	deleteUser.PersistentFlags().StringP("user-id", "u", "", k_userId)
	deleteUser.PersistentFlags().BoolP("json", "j", false, k_json)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error()+".")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "telstar-util",
	Short: "Utility program for interacting with the Telstar Server API. (c) John Newcombe 2024. Version: " + globals.Version,
	Long:  `Utility program for interacting with the Telstar Server API. (c) John Newcombe 2024.`,
}
