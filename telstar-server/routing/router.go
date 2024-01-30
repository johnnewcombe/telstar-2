package routing

import (
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar/session"
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/johnnewcombe/telstar-library/logger"
)

const ( // iota is reset to 0
	Undefined              = iota //0
	RouteMessageUpdated           // 1
	ValidPageRequest              // 3
	InvalidPageRequest            // 4
	InvalidCharacter              // 5
	BufferCharacterDeleted        //6
)

type PID struct {
	PageNumber int
	FrameId    string
}

type RouterRequest struct {
	InputByte        byte
	CurrentPageId    string
	HasFollowOnFrame bool // indicates that the current frame has a follow on frame
	RoutingTable     []int
	SessionId        string
}

type RouterResponse struct {
	Status        int
	RoutingBuffer string
	NewPageId     string
	HistoryPage   bool
	ImmediateMode bool
}

func ProcessRouting(request *RouterRequest, response *RouterResponse) error {

	var (
		err error
	)

	// main processing loop
	// there are two basic states, immediate mode should navigate immediately to the page based
	// on the routing table of the current page and buffer mode (non-immediate mode) where the user requests a
	// specific frame with e.g. *270#

	if !request.isValid() {
		return fmt.Errorf("routing request was invalid/empty")
	}
	// check for large data request e.g. http etc. and truncate
	truncateData(request, response)

	// check for asterisk, backspace etc. asterusk will switch to immediate mode
	// this does not add the byte to the buffer, it only handles special chars
	preProcessInputByte(request, response)
	if request.InputByte == globals.BS {
		// this is a special case so take an early bath
		return nil
	}

	// Start of routing process
	logger.LogInfo.Printf("Routing Start. [Current Page: %s [Message Buffer: %s] [Character to Process: %c]\r\n", request.CurrentPageId, response.RoutingBuffer, request.InputByte)

	if response.ImmediateMode {

		// gets new page Id from routing table of current frame
		if err = processImmediateMode(request, response); err != nil {
			return err
		}

		if response.Status == ValidPageRequest {
			// # has been received and the pageId is good
			logger.LogInfo.Printf("Frame Id is valid. [PageId: %s]", response.NewPageId)

		} else {

			// this only occurs if # has been entered and for some reason the pageId is invalid
			// this could occur if a command is entered e.g. *VT52# or *WEATHER# etc.

			logger.LogInfo.Printf("Frame Id is invalid. [PageId: %s]", response.NewPageId)

		}

	} else { // non-immediate mode (buffermode)

		if err = processBufferMode(request, response); err != nil {
			return err
		}
		// check some special page numbers
		// request to reload the current page
		if response.RoutingBuffer == "00" || response.RoutingBuffer == "09" {

			if err = selectReloadFrame(request, response); err != nil {
				return err
			}
			logger.LogInfo.Printf("Current page re-selected: [Page Id: %s]\r\n", response.NewPageId)

		}
		// check some special command sequences
		// back a page
		if response.RoutingBuffer == "" {

			// means that '*#' has been entered
			if err = selectPreviousFrame(request, response); err != nil {
				return err
			}
			logger.LogInfo.Printf("Previous frame selected (obtained from history): [Page Id: %s]", response.NewPageId)

		} else if strings.ToLower(response.RoutingBuffer) == "exit" { // means that '*EXIT#' has been enterd

			// code to execute commands goes here
			//
			// e.g.
			// selectMyCommand(&request, &response)
			logger.LogInfo.Printf("Command Exit was invoked.")
		}
	}
	// request to reload the current page
	//if response.RoutingBuffer == "*00" || response.RoutingBuffer == "*09" {

	//	selectReloadFrame(&request, &response)
	//	logger.LogInfo.Printf("Current page re-selected: [Page Id: %s]\r\n", response.newPageId)

	//} else if response.RoutingBuffer == "*\x5f" {

	//	selectPreviousFrame(&request, &response)
	//	logger.LogInfo.Printf("Previos frame selected (obtained from history): [Page Id: %s]", response.newPageId)
	/*
		} else if response.immediateMode && request.InputByte == 0x5f {

			selectFollowOnFrame(&request, &response)
			logger.LogInfo.Printf("Follow-On frame selected: [Page Id: %s]", response.newPageId)

		} else if response.immediateMode {

			processImmediateMode(&request, &response)

			if response.status == VALID_PAGE_REQUEST {
				// page id is good
				logger.LogInfo.Printf("[non-immediate mode] Frame Id is valid. [Buffer: %s]", response.RoutingBuffer)
			} else if response.status == INVALID_PAGE_REQUEST {
				logger.LogInfo.Printf("[non-immediate mode] Frame Id is invalid. [Buffer: %s]", response.RoutingBuffer)
			}

		} else {
			processNonImmediateMode(&request, &response)
			if response.status == VALID_PAGE_REQUEST {
				// page id is good
				logger.LogInfo.Printf("Frame Id is valid. [Buffer: %s]", response.RoutingBuffer)

			} else if response.status == INVALID_PAGE_REQUEST {
				// this only occurs if # has been entered and for some reason the pageId is invalid
				// this should never happen
				logger.LogInfo.Printf("Frame Id is invalid. [Buffer: %s]", response.RoutingBuffer)
			}
		}
	*/
	// routing complete
	logger.LogInfo.Printf("Routing finished [Message Buffer: %s]", response.RoutingBuffer)

	return nil
}

func ForceRoute(pageNumber int, frameId string, request *RouterRequest, response *RouterResponse) {

	// FIXME change frameId to a rune
	//  return err
	//  check for valid pageId

	// TODO check for err != nil
	response.RoutingBuffer, _ = utils.ConvertPidToPageId(pageNumber, frameId)

	response.trimBuffer()

	pageId := response.RoutingBuffer

	if utils.IsValidPageId(pageId) { // this will fail for commands like *vt52# etc

		response.NewPageId = pageId
		response.Status = ValidPageRequest

	} else {
		response.NewPageId = response.RoutingBuffer
		response.Status = InvalidPageRequest
	}

	// once complete reset to immediate mode
	response.ImmediateMode = true
}
func (response *RouterResponse) trimBuffer() {

	// Function to remove a leading * and trailing #
	response.RoutingBuffer = strings.TrimLeft(response.RoutingBuffer, "*")
	response.RoutingBuffer = strings.TrimRight(response.RoutingBuffer, "\x5f")

	return
}

func (response *RouterResponse) Clear() {
	//response.currentPageId = ""
	//response.followOnPageId = ""
	response.HistoryPage = false
	response.ImmediateMode = true
	response.NewPageId = ""
	response.RoutingBuffer = ""
	response.Status = 0

	return
}

func (request *RouterRequest) isValid() bool {

	// check routing table entries and input byte value
	if utils.IsValidPageId(request.CurrentPageId) &&
		len(request.RoutingTable) == 11 &&
		request.InputByte != 0 {

		//check routing table entries
		for _, i := range request.RoutingTable {

			s := fmt.Sprintf("%d", i)
			if len(s) > 9 {
				return false
			}
		}
		return true
	}
	return false

}

func truncateData(request *RouterRequest, response *RouterResponse) {

	// check for large data request e.g. http etc. and truncate
	if len(response.RoutingBuffer) > 11 {

		logger.LogInfo.Print("Message is too long, truncating.\r\n")

		// truncate the buffer
		//if response.immediateMode {
		response.RoutingBuffer = ""
		//} else {
		//	response.routingBuffer = "*" // asterisk is never placed in the buffer, so doesn't need to be here
		//}
	}
	return
}

func preProcessInputByte(request *RouterRequest, response *RouterResponse) {

	// this function handles special characters such as asterisk, backspace etc.
	// it doesn't add chars to the buffer, this is only done in buffer mode
	response.Status = Undefined

	if request.InputByte == 0x2a { // '*'

		// user looking for a particular page?
		response.ImmediateMode = false

		//start a new frame routing
		response.RoutingBuffer = ""

	} else if request.InputByte == 0x08 { // BS

		// process a backspace by removing a char from the routing buffer
		if len(response.RoutingBuffer) > 0 {
			response.RoutingBuffer = response.RoutingBuffer[:len(response.RoutingBuffer)-1]
			response.Status = BufferCharacterDeleted

			if len(response.RoutingBuffer) == 0{
				response.ImmediateMode = true
			}
		}
	}
	return
}

func selectReloadFrame(request *RouterRequest, response *RouterResponse) error {

	// current frame re-selected

	// reset to immediate mode
	response.ImmediateMode = true

	// re-display (00) / update (09) the current page
	// set page to be loaded to the current page id
	response.NewPageId = request.CurrentPageId
	response.Status = ValidPageRequest

	return nil
}

func selectPreviousFrame(request *RouterRequest, response *RouterResponse) error {

	// reset to immediate mode
	response.ImmediateMode = true

	// pop the most recent page
	pageId, ok := session.PopHistory(request.SessionId)

	// if the popped page is the current page then
	// pop again, note that history pages are not
	// put back into the history
	if pageId == request.CurrentPageId {
		pageId, ok = session.PopHistory(request.SessionId)
	}

	if ok {
		// indicate that this page id was retrieved from the history
		response.HistoryPage = true
		response.NewPageId = pageId

	} else {
		response.HistoryPage = false
		response.NewPageId = request.CurrentPageId
	}

	response.Status = ValidPageRequest

	return nil
}

func selectHashFrame(request *RouterRequest, response *RouterResponse) error {

	var (
		pageId string
		err    error
	)

	if pageId, err = getHashPageId(request); err != nil {
		return err
	}

	if utils.IsValidPageId(pageId) {

		response.NewPageId = pageId
		response.Status = ValidPageRequest

	} else {
		response.Status = InvalidPageRequest
	}

	return nil
}

func selectFollowOnFrame(request *RouterRequest, response *RouterResponse) error {

	var (
		pageId string
		err    error
	)
	// 'follow on frame' selected
	if pageId, err = GetFollowOnPageId(request.CurrentPageId); err != nil {
		return err
	}

	if utils.IsValidPageId(pageId) {

		response.NewPageId = pageId
		response.Status = ValidPageRequest

	} else {
		response.Status = InvalidPageRequest
	}

	//TODO is this needed?
	//response.trimBuffer()

	return nil
}

func processImmediateMode(request *RouterRequest, response *RouterResponse) error {

	var (
		err error
	)

	if !response.ImmediateMode {
		logger.LogError.Fatal("Calling processImmediateMode() when in Buffer Mode is invalid.")
	}

	if !request.isValid() {
		return fmt.Errorf("invalid routing request.")
	}

	if utils.IsNumeric(request.InputByte) {

		pageId := fmt.Sprintf("%da", request.RoutingTable[request.InputByte-0x30])

		//frame_info.page_id = self.current_frame.get_route_entry(ord(char) - 0x30)

		if utils.IsValidPageId(pageId) {

			response.NewPageId = pageId
			response.Status = ValidPageRequest

		} else {
			response.Status = InvalidPageRequest
		}
	} else if request.InputByte == 0x5f {

		if request.HasFollowOnFrame {
			if err = selectFollowOnFrame(request, response); err != nil {
				return err
			}
		} else {
			if err = selectHashFrame(request, response); err != nil {
				return err
			}
		}
	} else {
		// not found returned as the key entered was not a number 0-9
		response.Status = InvalidCharacter
	}
	return nil
}

func processBufferMode(request *RouterRequest, response *RouterResponse) error {

	if response.ImmediateMode {
		logger.LogError.Fatal("Calling processBufferMode() when in Immediate mode is invalid.")

	}
	if utils.IsAlphaNumeric(request.InputByte) {

		if request.InputByte == 0x5f {

			// remove any typed '*' and trailing '#'
			response.trimBuffer()
			pageId := response.RoutingBuffer + "a"

			if utils.IsValidPageId(pageId) { // this will fail for commands like *vt52# etc

				response.NewPageId = pageId
				response.Status = ValidPageRequest

			} else {
				response.NewPageId = response.RoutingBuffer
				response.Status = InvalidPageRequest
			}

			// once complete reset to immediate mode
			response.ImmediateMode = true

		} else {
			//all good so add inputChar to buffer
			response.RoutingBuffer += string(request.InputByte)

			// buffer mode *00 and *09 should update the new page ID to the current and set to Valid Page Request
			if response.RoutingBuffer == "*00" || response.RoutingBuffer == "*09" {

				// remove any typed '*' and trailing '#'
				response.trimBuffer()
				response.NewPageId = request.CurrentPageId
				response.Status = ValidPageRequest

				// once complete reset to immediate mode
				response.ImmediateMode = true

			} else {
				// still in non-immediate mode but waiting for the \x5f terminator
				response.Status = RouteMessageUpdated
			}
		}
	} else {
		response.Status = InvalidCharacter
	}
	return nil
}

func getHashPageId(request *RouterRequest) (string, error) {
	var (
		pageId string
	)
	pageId = fmt.Sprintf("%da", request.RoutingTable[10])
	if utils.IsValidPageId(pageId) {
		return pageId, nil
	}
	return "", errors.New("routing table entry for hash route is invalid")
}

func GetFollowOnPageId(pageId string) (string, error) {

	var (
		pageNumber int
		frameId    string
		err        error
	)

	/*
	   This function returns the frame that should be returned if the hash '#' (\x5f) key is pressed.
	   For example, if the current page is 200a, then 200b would be returned assuming the page
	   exists. If the page does not exist then the function will use the page specified in routing_table[10].
	*/

	if !utils.IsValidPageId(pageId) {
		return "", errors.New("invalid page id")
	}

	if pageNumber, frameId, err = utils.ConvertPageIdToPID(pageId); err != nil {
		return "", err
	}
	//frameId := []rune(pageId[len(pageId)-1:])
	frameIdAsc := int(frameId[0])
	pageNumberAsc := fmt.Sprintf("%d", pageNumber)

	//pageNumber := pageId[:len(pageId)-1]

	// update frame indicator and include zero page routing
	if frameIdAsc < 97 {
		return "", errors.New("invalid frame id")

	} else if frameIdAsc < 122 {
		frameIdAsc += 1
	} else {
		frameIdAsc = 97
		pageNumberAsc += "0"
	}
	return pageNumberAsc + string(rune(frameIdAsc)), nil
}
