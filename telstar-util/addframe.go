package main

/*
import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-util/http"
	"errors"
	"fmt"
)

func cmdAddFrame(apiUrl string, frameFile string, primary bool, token string) (ResponseData, error) {

	var (
		pid       types.Pid
		ok        bool
		err       error
		frameData string
		respData  ResponseData
		frame     types.Frame
		url       = apiUrl + "/frame"
	)

	if primary {
		url += "?db=primary"
	}

	if frameData, err = http.LoadText(frameFile); err != nil {
		return respData, err
	}

	if isEditTfFrame(frameData) {

		// we need the pid from the filename as edit.tf pages do not have a PID embedded within them
		// TODO Check the top line of the edit.tf for a PID? then we wouldn't need to rely on the filename
		//  or do both, i.e. no PID in edit.tf then check filename?
		if pid, ok = getPidFromFileName(frameFile); !ok {
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

	respData, err = http.Put(url, frameData, token)
	if err != nil {
		return err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}

	return nil
}


*/
