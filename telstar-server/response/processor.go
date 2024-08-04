package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-library/utils"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/renderer"
	"github.com/johnnewcombe/telstar/session"
	"net"
	"path/filepath"
	"time"
)

// TODO would like to dynamically fill a field before display.
//  An example would be the current date and time on the Welcome screen.
//  And on the Funds Transfer screen I want the Transaction Date to be filled with the next working day.
//  Later I would also like to fill a frame with values from a file.
//  An example is the screen that shows the latest transactions.

type Field struct {
	Value    string
	Complete bool
}
type ResponseData struct {
	FieldValues       []Field
	Complete          bool
	ActionOut         string
	currentFieldIndex int
	asterisk          bool
}

func (r *ResponseData) IncreaseFieldIndex() bool {
	if r.currentFieldIndex < len(r.FieldValues)-1 {
		r.currentFieldIndex++
		return true
	}
	return false
}

func (r *ResponseData) GetCurrentField() Field {
	return r.FieldValues[r.currentFieldIndex]
}

func (r *ResponseData) Clear() {
	r.FieldValues = []Field{}
	r.Complete = false
	r.currentFieldIndex = 0
}

func Process(sessionId string, conn net.Conn, inputByte byte, frame *types.Frame, responseData *ResponseData, baudRate int, settings config.Config) error {

	var (
		fieldCount           int
		currentValueLength   int
		currentResponseField types.ResponseField
		actionOut            string
		resultFrames         []types.Frame
		fieldFull            bool
		err                  error
	)

	fieldCount = len(frame.ResponseData.Fields)
	if fieldCount == 0 {
		return errors.New("No fields defined for this response page.")
	}

	if len(responseData.FieldValues) == 0 {
		// clear the response data for the next response frame
		// first time here for this page so initialsise the field value
		responseData.Clear()
		responseData.FieldValues = make([]Field, fieldCount, fieldCount)
	}

	currentResponseField = frame.ResponseData.Fields[responseData.currentFieldIndex]
	currentValueLength = len(responseData.FieldValues[responseData.currentFieldIndex].Value)

	// this should never be the case but belt and braces
	if responseData.currentFieldIndex >= fieldCount {
		return errors.New("The 'current field index' is invalid.")
	}

	if inputByte != globals.HASH {

		if inputByte == globals.BS {
			if currentValueLength > 0 { // backspace
				responseData.FieldValues[responseData.currentFieldIndex].Value =
					responseData.FieldValues[responseData.currentFieldIndex].Value[:currentValueLength-1]
				renderBuffer(conn, []byte{0x08, 0x20, 0x08}, baudRate)
			}
			return nil
		}

		// note if asterisk
		responseData.asterisk = inputByte == globals.ASTERISK

		if (currentResponseField.Type == "numeric" && utils.IsNumeric(inputByte)) ||
			(currentResponseField.Type == "alphanumeric" && utils.IsAlphaNumeric(inputByte)) {

			// add the value if there is room and then check to see if we are full
			if len(responseData.FieldValues[responseData.currentFieldIndex].Value) < currentResponseField.Length {
				responseData.FieldValues[responseData.currentFieldIndex].Value += string(inputByte)
				if currentResponseField.Password {
					renderBuffer(conn, []byte{globals.ASTERISK}, baudRate)
				} else {
					renderBuffer(conn, []byte{inputByte}, baudRate)
				}
			}
			if len(responseData.FieldValues[responseData.currentFieldIndex].Value) == currentResponseField.Length {
				fieldFull = true
			}
		}
	}

	// check if field is complete, move to next field if 0x5f or if length
	// of field reached AND autosubmit is true
	if inputByte == globals.HASH || (fieldFull && currentResponseField.AutoSubmit) {

		// if cancel sequence is detected by checking the asterisk flag, this will only be set
		// if the previous char was an asterisk. This will be the case no matter what type of
		// field this is.
		if responseData.asterisk {

			// by copying the cancel pid to the action pid, the cancel pid will get
			// rendered whem we exit this function
			frame.ResponseData.Action.PostActionFrame = frame.ResponseData.Action.PostCancelFrame

			// nothing more to do so mark as such
			responseData.Complete = true
			return nil
		}

		// field isn't complete if it is a required field and is empty
		if !(currentResponseField.Required && currentValueLength == 0) {

			responseData.FieldValues[responseData.currentFieldIndex].Complete = true

			// is the response complete?
			if !responseData.IncreaseFieldIndex() {

				//that was the last field so we are complete
				args := append(frame.ResponseData.Action.Args)
				//args := append(frame.ResponseData.Action.Args, responseData.FieldValues...)

				// carry out the action
				for _, v := range responseData.FieldValues {
					args = append(args, v.Value)
				}
				// any previous frames need to be removed?
				session.ClearCache(sessionId)
				if actionOut, err = action(sessionId, frame, args, settings); err != nil {
					file, _ := filepath.Abs(frame.ResponseData.Action.Exec)
					return fmt.Errorf("%s, %s, %v", file, actionOut, err)
				}

				if err = json.Unmarshal([]byte(actionOut), &resultFrames); err != nil {
					// we will get this error if the search term wasn't found
					var resultFrame types.Frame
					if err = json.Unmarshal([]byte(actionOut), &resultFrame); err != nil {

						logger.LogInfo.Println("the response frame action did not return frame(s)")

						// some commands such as internal commands (login, mail etc)
						// may not return page data so just continue

					} else {
						resultFrames = []types.Frame{resultFrame}
					}
				}

				// we are complete so place any resulting frame into session cache
				// remember these are temporary pages, these are typically follow-on frames
				// from the next frame to be displayed i.e. the response frames post action frame is rendered and
				// Note that some commands may not return anything at all or some other data type which
				// is effectively ignored.

				for f := 0; f < len(resultFrames); f++ {
					frame := resultFrames[f]
					session.AddFrameToCache(sessionId, frame)
				}

				// complete so mark as such
				responseData.Complete = true
				return nil

			} else {

				// move cursor to next field
				renderer.PositionCursor(conn, frame.ResponseData.Fields[responseData.currentFieldIndex].HPos,
					frame.ResponseData.Fields[responseData.currentFieldIndex].VPos, !settings.Server.DisableVerticalRollOver)
			}
		}
	}

	return nil
}

func renderBuffer(conn net.Conn, buffer []byte, baudRate int) error {

	for _, b := range buffer {
		if _, err := conn.Write([]byte{b}); err != nil {
			return err
		}
		// slow down to match the baud rate
		time.Sleep(time.Duration(baudRate))
	}
	return nil
}
