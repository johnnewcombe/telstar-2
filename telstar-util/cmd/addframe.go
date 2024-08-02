package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var addFrame = &cobra.Command{
	Use:   "add-frame",
	Short: "Adds a single frame to the currently logged in system.",
	Long: `
Adds a single frame to the currently logged in system. See the login command.
If the frame already exists, the frame will be updated, if it does not exist 
then it will be created.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl        string
			primary       bool
			includeUnsafe bool
			token         string
			frameData     string
			err           error
			respData      network.ResponseData
			frame         types.Frame
			filename      string
		)

		if filename, err = cmd.Flags().GetString("filename"); err != nil {
			return err
		}

		//load token - don't want to report errors here as we want an unauthorised status code to be returned
		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}
		if includeUnsafe, err = cmd.Flags().GetBool("include-unsafe"); err != nil {
			return err
		}
		if primary, err = cmd.Flags().GetBool("primary"); err != nil {
			return err
		}

		apiUrl = apiUrl + "/frame"
		if primary {
			apiUrl += "?db=primary"
		}

		if frameData, err = loadText(filename); err != nil {
			return err
		}

		if isEditTfFrame(frameData) {
			return errors.New("edit.t frames are only supported as part of a standard json frame definition")
		}

		// validate the frameData
		if frame, err = parseFrame(frameData); err != nil {
			err = fmt.Errorf("invalid frame data: %v", err)
			return err
		}

		// note that setting page ID to "" prevents any restrictions in updating
		if respData, err = addSingleFrameJson(apiUrl, frame, includeUnsafe, token); err != nil {
			return err
		}

		stdOut(cmd, respData, nil)
		return nil
	},
}

// addSingleFrameJson This function accepts pageId, which if set will restrict updates to that specific page
func addSingleFrameJson(apiUrl string, frame types.Frame, includeUnSafe bool, token string) (network.ResponseData, error) {

	var (
		err       error
		respData  network.ResponseData
		frameData []byte
		//ok       bool
	)

	if !frame.IsValid() {
		err = errors.New("invalid frameData")
		return respData, err
	}

	if !includeUnSafe {

		// can't bulk upload these as users could use their frames to run system utils etc.
		if strings.ToLower(frame.FrameType) == "response" {
			err = fmt.Errorf("frame %d%s is a response frame, set option --include-unsafe to force", frame.PID.PageNumber, frame.PID.FrameId)
			return respData, err
		}

		// can't bulk upload these as users could use their frames to run system utils etc.
		if strings.ToLower(frame.FrameType) == "gateway" {
			err = fmt.Errorf("frame %d%s is a gateway frame, set option --include-unsafe to force", frame.PID.PageNumber, frame.PID.FrameId)
			return respData, err
		}
	}

	if frameData, err = frame.Dump(); err != nil {
		return respData, err
	}

	respData, err = network.Put(apiUrl, string(frameData), token)
	if err != nil {
		return respData, err
	}

	return respData, nil

}
