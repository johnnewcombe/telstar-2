package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
)

var getFrames = &cobra.Command{
	Use:   "get-frames",
	Short: "Returns multiple frames from the currently logged in system.",
	Long: `
Returns multiple frame from the currently logged in system. See the login command.

A destination directory is specified as a location for the frames to be stored.
Frames are returned in json format and each frame is stored as a single json file
with the filename in <page number><frame id>.json

These files can be specified as a source to the the add-frame command and the 
directory can be specified as a source to the add-frames command
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl      string
			respData    network.ResponseData
			jsonData    []map[string]interface{}
			data        []byte
			count       int
			err         error
			token       string
			destination string
			primary     bool
		)

		//load token - don't want to report errors here as we want an unauthorised status code to be returned
		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if destination, err = cmd.Flags().GetString("destination"); err != nil {
			return err
		}
		if primary, err = cmd.Flags().GetBool("primary"); err != nil {
			return err
		}

		apiUrl += "/frame"
		if primary {
			apiUrl += "?db=primary"
		}

		respData, err = network.Get(apiUrl, token)
		if err != nil {
			return err
		}

		if respData.StatusCode < 200 || respData.StatusCode > 299 {
			// Render HTTP response and exit
			stdOut(cmd, respData, nil)
			return nil
		}

		//parse the json array of tmp as unstructured data
		if err = json.Unmarshal([]byte(respData.Body), &jsonData); err != nil {
			return err
		}

		for _, frame := range jsonData {

			// get the pid so that we can create the filename
			pid := frame["pid"].(map[string]interface{}) // this is a type assertion to convert from interface{}

			// get underlying data types
			pageNo := int(pid["page-no"].(float64))
			frameId := pid["frame-id"].(string)

			// convert each frame back to json so that each frame can be saved individually
			if data, err = json.MarshalIndent(frame, "", "    "); err != nil {
				return err
			}

			// create filename
			filename := filepath.Join(destination, fmt.Sprintf("%d%s.json", pageNo, frameId))

			// save the file
			if err = saveText(filename, string(data)); err != nil {
				return err
			}
			count++
		}

		result := map[string]string{
			"Saved": strconv.Itoa(count),
		}
		stdOut(cmd, respData, result)
		return nil
	},
}
