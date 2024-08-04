package renderer

import (
	"bitbucket.org/johnnewcombe/telstar-library/convert"
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/text"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar/apps"
	"bitbucket.org/johnnewcombe/telstar/config"
	"bitbucket.org/johnnewcombe/telstar/session"
	"bitbucket.org/johnnewcombe/telstar/synchronisation"
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
)

type RenderOptions struct {
	HasFollowOnFrame bool
	BaudRate         int
}

// Render
func Render(ctx context.Context, conn net.Conn, wg *synchronisation.WaitGroupWithCount, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) {

	defer wg.Done()
	if utils.IsValidPageId(frame.GetPageId()) {

		renderHeader(ctx, conn, frame, sessionId, settings, options)
		renderTitle(ctx, conn, frame, sessionId, settings, options)
		renderContent(ctx, conn, frame, sessionId, settings, options)
		renderFooter(ctx, conn, frame, sessionId, settings, options)

		if frame.FrameType != globals.FRAME_TYPE_TEST && frame.FrameType != globals.FRAME_TYPE_RESPONSE {
			wg.Add(1)
			go RenderSystemMessage(ctx, conn, wg, frame.NavMessage, settings, options)
		}

		if frame.FrameType == "response" && len(frame.ResponseData.Fields) > 0 {
			if err := PositionCursor(conn, frame.ResponseData.Fields[0].HPos, frame.ResponseData.Fields[0].VPos, !settings.Server.DisableVerticalRollOver); err != nil {
				logger.LogError.Print(err)
			}
		}

	} else {
		// page missing so just Render navigation to display 'not found'
		//renderNavigationMessage(ctx, conn, frame, settings, user, options)
		logger.LogError.Print(errors.New("render requested for an invalid frame"))
		//wg.Add(1)
		//RenderTransientSystemMessage(ctx, conn, wg, settings.Strings.NavMessageNotFound, settings.Strings.NavMessageSelectPage, options)
	}

	return
}

// RenderSystemMessage
func RenderSystemMessage(ctx context.Context, conn net.Conn, wg *synchronisation.WaitGroupWithCount, message string, settings config.Config, options RenderOptions) {
	// FIXME With merged pages that are 960 char long, rendering this causes a scroll, this shouldn't happen if HOME/VTAB is used, should it?
	var (
		err         error
		cursorChars strings.Builder
		lastChar    string
	)

	logger.LogInfo.Printf("System Message: %s\r\n", message)

	defer wg.Done()

	// position the cursor to the bottom row, column 0
	if err = PositionCursor(conn, 0, globals.ROWS-1, !settings.Server.DisableVerticalRollOver); err != nil {
		logger.LogError.Print(err)
		return
	}

	// convert the message markup and swap the pound signs
	if message, err = convert.MarkupToRawV(message); err != nil {
		logger.LogError.Print(err)
		return
	}
	message = strings.ReplaceAll(message, string(globals.POUND), string(globals.HASH))
	cursorPos := text.GetDisplayLen(message)

	// all navigation messages start with a cursor-off, the user may define
	// a cursor-on for a specific nav message (see settings)
	// if last char of message is cursor on, remove it, once all of the
	// cursor positioning has been done, pop it back (see below)
	if message[len(message)-1:] == string(globals.CURON) {
		message = message[:len(message)-1]
		lastChar = string(globals.CURON)
	}

	// pad line to 39 chars this ensures we cover any previously displayed message
	// note that we have to pad, before setting the cursor position otherwise
	message = text.PadTextRight(message, globals.COLS-1)

	// add to the string builder
	cursorChars.WriteByte(globals.CUROFF)
	cursorChars.WriteString(message)

	// put the cursor at the beginning of the line and HT to correct position
	// most nav messages are < 20 chars so this is generally the most efficient
	cursorChars.WriteByte(globals.CR)
	for c := 0; c < cursorPos; c++ {
		cursorChars.WriteString(string(globals.HT))
	}

	// put back the cursor on if it was set.
	cursorChars.WriteString(lastChar)
	renderBuffer(ctx, conn, []byte(cursorChars.String()), settings, options)

}

// RenderTransientSystemMessage displays the specified message that is then replaced
// a second later with the specified follow-on message
func RenderTransientSystemMessage(ctx context.Context, conn net.Conn, wg *synchronisation.WaitGroupWithCount, message string, followOnMessage string, settings config.Config, options RenderOptions) {

	defer wg.Done()

	if len(message) > 0 {
		wg.Add(1)
		RenderSystemMessage(ctx, conn, wg, message, settings, options)
	}

	if len(followOnMessage) > 0 {
		// if the RenderSystemMessage call (above) was cancelled we will
		// have a ctx.Err so make make a tidy exit
		if ctx.Err() != nil {
			return
		}

		// all good so complete the system message
		time.Sleep(1000 * time.Millisecond)
		wg.Add(1)
		RenderSystemMessage(ctx, conn, wg, followOnMessage, settings, options)
	}
}

func renderHeader(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) {

	var (
		header     string
		headerText string
		metadata   string
		cls        string
		pageId     string
		cost       string
		err        error
	)
	defer ctx.Done()

	// test pages just send home/clear and make a quick exit
	if frame.FrameType == globals.FRAME_TYPE_TEST {

		if !frame.DisableClear {
			cls = string(globals.CLS)
		} else {
			cls = string(globals.HOME)
		}

		buffer := []byte(cls)
		renderBuffer(ctx, conn, buffer, settings, options)

		return
	}

	//actual header text
	if len(frame.HeaderText) > 0 {
		if headerText, err = convert.MarkupToRawV(frame.HeaderText); err != nil {
			logger.LogError.Print(err)
		}
	} else {
		if headerText, err = convert.MarkupToRawV(settings.Server.Strings.DefaultHeaderText); err != nil {
			logger.LogError.Print(err)
		}
	}

	if !frame.DisableClear && frame.FrameType == globals.FRAME_TYPE_INITIAL {
		//	if options.ClearScreen && frame.FrameType == "initial" {
		metadata = getMetaData(time.Now())
		cls = string(globals.CLS) + metadata + string(globals.CLS)
	} else if !frame.DisableClear {
		cls = string(globals.CLS)
	} else {
		cls = string(globals.HOME)
	}

	if !settings.Server.HidePageId {
		pageId = frame.GetPageId()
	}
	pageId = text.PadTextRight(pageId, 11)

	if !settings.Server.HideCost {
		cost = fmt.Sprintf("%dp", frame.Cost)
	}

	cost = text.PadTextLeft(cost, 4)

	// header is formatted such that we have a 10 char pageId a spc and a 3 char cost (e.g. 14 chars)
	// line length is 40, therefore header needs to be truncated or padded to 40 -14 == 26 but this doesn't
	// include non printing chars e.g.ctrl chars
	headerText = text.PadTextRight(headerText, 24)

	header = fmt.Sprintf("%s%s%s%s%s ", cls, headerText, string(globals.CUROFF), pageId, cost)

	// TODO pad and add page number and add pad to 40 chars HEADER shouldn't be shorter as it may be used without a CLS
	// and we need to make sure the header replaces any previous one completely.

	// no page number on initial frame or when page number is switched off
	// could use home and htab of fit it into the header

	buffer := []byte(header)
	//fmt.Printf("%v\r\n", buffer)

	renderBuffer(ctx, conn, buffer, settings, options)

}

func renderTitle(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) {

	var (
		data string
		err  error
	)

	// no titles for test pages
	if frame.FrameType == globals.FRAME_TYPE_TEST {
		return
	}

	// split the type as this allows for comma separated params.
	dataType, rows := utils.ParseDataType(frame.Title.Type)

	// Render the content
	switch dataType {
	case "markup":
		if data, err = convert.MarkupToRawV(frame.Title.Data); err != nil {
			logger.LogError.Print(err)
			return
		}
		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return
	case globals.CONTENT_TYPE_RAW, globals.CONTENT_TYPE_RAWV:

		data = frame.Title.Data

		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return

	case globals.CONTENT_TYPE_RAWT:
		if data, err = convert.RawTToRawV(frame.Title.Data, 0, 23, 0, 39, true); err != nil {
			logger.LogError.Print(err)
			return
		}
		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return

	case globals.CONTENT_TYPE_EDITTF, "edittf", globals.CONTENT_TYPE_ZXNET:
		// for an edit.tf title field a number can be added to the type e.g. edit.tf,7 which will take the first 7
		// rows rather than the default (usually 4). This can be used to have odd frames with a larger number of rows
		// used for the title.
		if rows == 0 {
			rows = settings.Server.EditTfTitleRows
		}
		// edit.tf mode for title only returns rows 1 to 4 inclusive from the edit.tf frame.
		if data, err = convert.EdittfToRawV(frame.Title.Data, 1, rows, !settings.Server.Antiope); err != nil {
			logger.LogError.Print(err)
			return
		}
		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
	}

	return
}

func renderContent(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) {

	var (
		data   string
		rowEnd int
		err    error
	)

	// Render the content
	switch frame.Content.Type {
	case "markup":
		// if markup then convert to raw
		if data, err = convert.MarkupToRawV(frame.Content.Data); err != nil {
			logger.LogError.Print(err)
			return
		}
		// swap placeholders
		data = populatePlaceholders(data, settings, sessionId, options)
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return
	case globals.CONTENT_TYPE_RAW, globals.CONTENT_TYPE_RAWV:
		data = populatePlaceholders(data, settings, sessionId, options)
		renderBuffer(ctx, conn, []byte(frame.Content.Data), settings, options)
		return
	case "rawT":
		if data, err = convert.RawTToRawV(frame.Content.Data, 0, 23, 0, 39, true); err != nil {
			logger.LogError.Print(err)
			return
		}
		data = populatePlaceholders(data, settings, sessionId, options)
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return

	case globals.CONTENT_TYPE_EDITTF, "edittf", globals.CONTENT_TYPE_ZXNET:
		if frame.FrameType == globals.FRAME_TYPE_TEST {
			rowEnd = 23
		} else {
			rowEnd = 22
		}
		// edit.tf is teletext, so get full teletext page
		if data, err = convert.EdittfToRawV(frame.Content.Data, 1, rowEnd, !settings.Server.Antiope); err != nil {
			logger.LogError.Print(err)
			return
		}
		data = populatePlaceholders(data, settings, sessionId, options)
		renderBuffer(ctx, conn, []byte(data), settings, options)
	}
	return
}

func renderFooter(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) {

	var (
		data string
		err  error
	)

	// split the type as this allows for comma separated params.
	dataType, rows := utils.ParseDataType(frame.Footer.Type)

	// Render the content
	switch dataType {
	case "markup":
		if data, err = convert.MarkupToRawV(frame.Footer.Data); err != nil {
			logger.LogError.Print(err)
			return
		}
		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return
	case globals.CONTENT_TYPE_RAW, globals.CONTENT_TYPE_RAWV:

		data = frame.Footer.Data

		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return

	case "rawT":
		if data, err = convert.RawTToRawV(frame.Content.Data, 0, 23, 0, 39, true); err != nil {
			logger.LogError.Print(err)
			return
		}
		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
		return

	case globals.CONTENT_TYPE_EDITTF, "edittf", globals.CONTENT_TYPE_ZXNET:
		// for an edit.tf title field a number can be added to the type e.g. edit.tf,7 which will take the first 7
		// rows rather than the default (usually 4). This can be used to have odd frames with a larger number of rows
		// used for the title.
		if rows == 0 {
			// FIXME we need a separate setting for this
			rows = settings.Server.EditTfTitleRows
		}
		// edit.tf mode for title only returns rows 1 to 4 inclusive from the edit.tf frame.
		if data, err = convert.EdittfToRawV(frame.Footer.Data, 1, rows, !settings.Server.Antiope); err != nil {
			logger.LogError.Print(err)
			return
		}
		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				logger.LogError.Print(err)
			}
		}
		renderBuffer(ctx, conn, []byte(data), settings, options)
	}

	return
}

func renderBuffer(ctx context.Context, conn net.Conn, buffer []byte, settings config.Config, options RenderOptions) {

	var (
		err error
	)

	if settings.Server.Antiope {
		if buffer, err = convert.RawVToAntiope(buffer); err != nil {
			logger.LogError.Print(err)
			return
		}
	}

	for _, b := range buffer {

		// process any requested cancellation
		select {
		case <-ctx.Done():
			// channel has a true, so cancel
			logger.LogInfo.Print("Rendering cancelled.")
			return // channel closed so cancel
		default:
		}

		if _, err := conn.Write([]byte{b}); err != nil {
			logger.LogError.Print(err)
		}

		// slow down to match the baud rate
		time.Sleep(time.Duration(options.BaudRate))

	}

	return
}

/*
func MergeRawVold(rawVData string, mergeData []string, rows int) (string, error) {

	var (
		err      error
		data     string
		layers   []string
		rawTData string
	)

	// apply any merge-data
	if mergeData != nil {

		// convert main data field to RawT
		if rawTData, err = convert.RawVToRawT(rawVData); err != nil {
			return "", err
		}
		// do the same for each of the merge fields
		for i := 0; i < len(mergeData); i++ {

			// convert all the merge layers to rawT via rawV
			if mergeData[i], err = convert.MarkupToRawV(mergeData[i]); err != nil {
				return "", err
			}
			if mergeData[i], err = convert.RawVToRawT(mergeData[i]); err != nil {
				return "", err
			}
		}

		// combine all of the layers
		layers = append(layers, rawTData)
		layers = append(layers, mergeData...)

		if data, err = convert.RawTMerge(layers...); err != nil {
			return "", err
		}
		// convert main data field back to RawV
		if data, err = convert.RawTToRawV(data, 0, rows-1, 0, 39, true); err != nil {
			return "", err
		}
	}
	return data, nil
}
*/

func PositionCursor(conn net.Conn, x int, y int, useRollover bool) error {

	if _, err := conn.Write([]byte{globals.HOME}); err != nil {
		return err
	}

	if useRollover {
		// TODO: Add rollover support based on the useRollover parameter
		// do some rolover calcs etc
		// e.g. if y > 12 use globals.VTAB * n
	} else {
		for vpos := 0; vpos < y; vpos++ {
			if _, err := conn.Write([]byte{globals.LF}); err != nil {
				return err
			}
		}
	}

	if useRollover {
		// TODO: Add rollover support based on the useRollover parameter
		// do some rollover calcs etc
		// e.g. if x > 20 use global.BS x n + globals.LF
	} else {
		for hpos := 0; hpos < x; hpos++ {
			if _, err := conn.Write([]byte{globals.HT}); err != nil {
				return err
			}
		}
	}
	if _, err := conn.Write([]byte{globals.CURON}); err != nil {
		return err
	}
	return nil
}

func populatePlaceholders(data string, settings config.Config, sessionId string, options RenderOptions) string {

	// always use the magical reference date for layouts
	// Mon Jan 2 15:04:05 MST 2006
	now := time.Now()
	user := session.GetCurrentUser(sessionId)

	data = strings.ReplaceAll(data, "[SERVER]", settings.Server.DisplayName)
	data = strings.ReplaceAll(data, "[DATE]", now.Format("2 Jan 2006"))
	data = strings.ReplaceAll(data, "[TIME]", now.Format("15:04"))
	data = strings.ReplaceAll(data, "[GREETING]", getGreeting(time.Now()))
	data = strings.ReplaceAll(data, "[NAME]", user.Name) //FIXME add user

	if strings.Contains(data, "[SYSINFO]") {
		data = strings.ReplaceAll(data, "[SYSINFO]", getSysInfo(settings, user, options))
	}

	// also app specific placeholders
	if strings.Contains(data, "[SHOP.PURCHASES]") {
		data = strings.ReplaceAll(data, "[SHOP.PURCHASES]", apps.ShopGetPurchases(sessionId, settings))
	}

	return data
}

func getGreeting(t time.Time) string {

	if t.Hour() < 12 {
		return "GOOD MORNING"
	} else if t.Hour() < 17 {
		return "GOOD AFTERNOON"
	} else {
		return "GOOD EVENING"
	}
}

func getMetaData(t time.Time) string {

	// example '20060102T1504Z'  taken at 15.40 GMT on 20 January 2006
	// always use the magical reference date for layouts
	// Mon Jan 2 15:04:05 MST 2006
	return t.Format("20060102T1504Z")
}

func getSysInfo(settings config.Config, user types.User, options RenderOptions) string {

	var (
		sb   strings.Builder
		baud string
	)
	ver, err := globals.GetVersion()
	if err != nil {
		logger.LogError.Printf("error loading version file %v", err)
		ver = ""
	}

	if options.BaudRate > 0 && options.BaudRate < 3000000 {
		baud = "4800"
	} else if options.BaudRate > 3000000 && options.BaudRate < 6000000 {
		baud = "2400"
	} else if options.BaudRate > 6000000 {
		baud = "1200"
	} else {
		baud = "MAX"
	}

	sb.WriteString(fmt.Sprintf("        Version : %s-%s\r\n", runtime.GOARCH, ver))
	sb.WriteString(fmt.Sprintf("         Server : %s\r\n", settings.Server.DisplayName))
	sb.WriteString(fmt.Sprintf("      Baud Rate : %s\r\n", baud))
	sb.WriteString(fmt.Sprintf("        User ID : %s\r\n", user.UserId))
	sb.WriteString(fmt.Sprintf("      User Name : %s\r\n", user.Name))
	sb.WriteString(fmt.Sprintf("      Base Page : %d\r\n", user.BasePage))
	sb.WriteString(fmt.Sprintf("       Database : %s\r\n", strings.ToUpper(settings.Database.Collection)))

	//return fmt.Sprintf("CPU Usage: %f%%\r\n   Busy: %f\r\n   Total: %f\n", cpuUsage, totalTicks-idleTicks, totalTicks)
	return sb.String()
}
