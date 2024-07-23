package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var addFrame = &cobra.Command{
	Use:   "add-frame",
	Short: "Adds a single frame to the currently logged in system.",
	Long: `
Adds a single frame to the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			err      error
			respData network.ResponseData
			filename string
		)

		if filename, err = cmd.Flags().GetString("filename"); err != nil {
			return err
		}
		if respData, err = addSingleFrame(cmd, filename); err != nil {
			return err
		}

		fmt.Printf(globals.Response, respData.Status)
		return nil
	},
}

func addSingleFrame(cmd *cobra.Command, filename string) (network.ResponseData, error) {

	var (
		apiUrl    string
		primary   bool
		token     string
		err       error
		respData  network.ResponseData
		frameData string
		frame     types.Frame
		pid       types.Pid
		ok        bool
	)

	//load token - don't want to report errors here as we want an unauthorised status code to be returned
	token, _ = loadText(globals.TOKENFILE)

	if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
		return respData, err
	}

	if primary, err = cmd.Flags().GetBool("primary"); err != nil {
		return respData, err
	}

	apiUrl = apiUrl + "/frame"
	if primary {
		apiUrl += "?db=primary"
	}

	if frameData, err = loadText(filename); err != nil {
		return respData, err
	}

	if isEditTfFrame(frameData) {

		// we need the pid from the filename as edit.tf pages do not have a PID embedded within them
		// TODO Check the top line of the edit.tf for a PID? then we wouldn't need to rely on the filename
		//  or do both, i.e. no PID in edit.tf then check filename?
		if pid, ok = GetPidFromFileName(filename); !ok {
			return respData, errors.New("filename format error")
		}

		if frameData, err = createEditTfFrame(pid, frameData); !ok {
			return respData, err
		}
	}

	// validate the frameData
	if frame, err = parseFrame(frameData); err != nil {
		err = fmt.Errorf("invalid frameData: %v", err)
		return respData, err
	}
	if !frame.IsValid() {
		err = errors.New("invalid frameData")
		return respData, err
	}

	respData, err = network.Put(apiUrl, frameData, token)
	if err != nil {
		return respData, err
	}

	return respData, nil
}
