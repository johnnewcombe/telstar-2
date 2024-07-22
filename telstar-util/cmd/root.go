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
	k_password    = "Password."
	k_pageId      = "Page ID."
	k_primary     = "Specify the primary database."
	k_source      = "Directory or GIT repository containing the json files to upload."
	k_destination = "Directory where the json files should be stored."
	k_fileName    = "Filename of json file."
	k_json        = "Output results in json"
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
	//addFrames.PersistentFlags().String("url", "", k_apiUrl)
	addFrames.PersistentFlags().StringP("source", "s", "", k_source)
	addFrames.PersistentFlags().Bool("primary", false, k_primary)

	// Add User
	//addUser.PersistentFlags().String("url", "", k_apiUrl)
	addUser.PersistentFlags().StringP("user-id", "u", "", k_userId)
	addUser.PersistentFlags().StringP("password", "p", "", k_password)

	// Delete Frame
	//deleteFrame.PersistentFlags().String("url", "", k_apiUrl)
	deleteFrame.PersistentFlags().String("page-id", "", k_pageId)
	deleteFrame.PersistentFlags().Bool("primary", false, k_primary)

	// Delete User
	//deleteUser.PersistentFlags().String("url", "", k_apiUrl)
	deleteUser.PersistentFlags().StringP("user-id", "u", "", k_userId)

	// Get Frame <url> <page id> [primary|secondary]"
	//getFrame.PersistentFlags().String("url", "", k_apiUrl)
	getFrame.PersistentFlags().String("page-id", "", k_pageId)
	getFrame.PersistentFlags().Bool("primary", false, k_primary)

	// Get Frames
	//getFrames.PersistentFlags().String("url", "", k_apiUrl)
	getFrames.PersistentFlags().StringP("destination", "d", "", k_destination)
	getFrames.PersistentFlags().Bool("primary", false, k_primary)

	// Get Status
	//getStatus.PersistentFlags().String("url", "", k_apiUrl)

	// Login
	//login.PersistentFlags().String("url", "", k_apiUrl)
	login.PersistentFlags().StringP("user-id", "u", "", k_userId)
	login.PersistentFlags().StringP("password", "p", "", k_password)

	// Purge Frames
	//purgeFrames.PersistentFlags().String("url", "", k_apiUrl)
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
