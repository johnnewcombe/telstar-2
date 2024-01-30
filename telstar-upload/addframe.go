package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"errors"
	"fmt"
	"strings"
)

var framesAdded = make(map[string]bool)

func cmdAddFrame(apiUrl string, frameFile string, primary bool, token string, basePage int, includeUnsafe bool) (ResponseData, error) {

	var (
		pid       Pid
		ok        bool
		err       error
		pageId    string
		frameData string
		respData  ResponseData
		frame     Frame
		url       = apiUrl + "/frame"
	)

	if primary {
		url += "?db=primary"
	}

	if frameData, err = loadText(frameFile); err != nil {
		return respData, err
	}

	if isEditTfFrame(frameData) || isZxNetFrame(frameData) {

		// we need the pid from the filename as edit.tf pages do not have a PID embedded within them
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

	if !utils.PageInScope(basePage, frame.PID.PageNumber) {
		return respData, fmt.Errorf("page number %d is not in scope, ignored", frame.PID.PageNumber)
	}

	if !frame.IsValid() {
		err = fmt.Errorf("invalid frame data for frame %d%s, ignored", frame.PID.PageNumber, frame.PID.FrameId)
		return respData, err
	}

	if !includeUnsafe {

		// can't bulk upload these as users could use their frames to run system utils etc.
		if strings.ToLower(frame.FrameType) == "response" {
			err = fmt.Errorf("frame %d%s is a response frame, ignored", frame.PID.PageNumber, frame.PID.FrameId)
			return respData, err
		}

		// can't bulk upload these as users could use their frames to run system utils etc.
		if strings.ToLower(frame.FrameType) == "gateway" {
			err = fmt.Errorf("frame %d%s is a gateway frame, ignored", frame.PID.PageNumber, frame.PID.FrameId)
			return respData, err
		}
	}

	if frame.PID.PageNumber >= 100 {
		if pageId, err = utils.ConvertPidToPageId(frame.PID.PageNumber, frame.PID.FrameId); err != nil {
			return respData, err
		}
		framesAdded[pageId] = true // this is a fudge but it gives us a map that returns false if the requested item ddoes not exist
	}

	respData, err = put(url, frameData, token)
	if err != nil {
		return respData, err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}

	logger.LogInfo.Printf("Adding frame %d%s.", frame.PID.PageNumber, frame.PID.FrameId )

	return respData, nil
}
