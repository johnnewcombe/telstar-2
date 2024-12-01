package renderer

import (
	"context"
	"errors"
	"fmt"
	"github.com/johnnewcombe/telstar-library/convert"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/text"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-library/utils"
	"github.com/johnnewcombe/telstar/apps"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/session"
	"github.com/johnnewcombe/telstar/synchronisation"
	"net"
	"runtime"
	"strings"
	"time"
)

type RenderOptions struct {
	HasFollowOnFrame bool
	BaudRate         int
}

type RenderResults []error

func Render(ctx context.Context, conn net.Conn, wg *synchronisation.WaitGroupWithCount, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions, chResult chan<- RenderResults) {

	var (
		err           error
		renderResults []error
	)

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "Render")
	}

	renderResults = make([]error, 0)

	defer wg.Done()

	if utils.IsValidPageId(frame.GetPageId()) {

		if err = renderHeader(ctx, conn, frame, sessionId, settings, options); err != nil {
			renderResults = append(renderResults, err)
		}

		if err = renderTitle(ctx, conn, frame, sessionId, settings, options); err != nil {
			renderResults = append(renderResults, err)
		}
		if err = renderContent(ctx, conn, frame, sessionId, settings, options); err != nil {
			renderResults = append(renderResults, err)
		}
		if err = renderFooter(ctx, conn, frame, sessionId, settings, options); err != nil {
			renderResults = append(renderResults, err)
		}

		if frame.FrameType != globals.FRAME_TYPE_TEST && frame.FrameType != globals.FRAME_TYPE_RESPONSE {
			//wg.Add(1)
			renderSystemMessage(ctx, conn, frame.NavMessage, sessionId, settings, options)
		}

		if frame.FrameType == "response" && len(frame.ResponseData.Fields) > 0 {
			if err = PositionCursor(conn, frame.ResponseData.Fields[0].HPos, frame.ResponseData.Fields[0].VPos, !settings.Server.DisableVerticalRollOver); err != nil {
				renderResults = append(renderResults, err)
			}
		}

	} else {
		renderResults = append(renderResults, errors.New("render requested for an invalid frame"))
	}

	// have some errors so send them back to listener.go
	if len(renderResults) > 0 {
		chResult <- renderResults
	}

	return
}

// RenderTransientSystemMessage displays the specified message that is then replaced
// a second later with the specified follow-on message
func RenderTransientSystemMessage(ctx context.Context, conn net.Conn, wg *synchronisation.WaitGroupWithCount, message string, followOnMessage string, sessionId string, settings config.Config, options RenderOptions, chResult chan<- RenderResults) {

	var (
		err           error
		renderResults []error
	)

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "RenderTransientSystemMessage")
	}

	renderResults = make([]error, 0)

	defer wg.Done()

	if len(message) > 0 {
		if err = renderSystemMessage(ctx, conn, message, sessionId, settings, options); err != nil {
			renderResults = append(renderResults, err)
		}
	}

	if len(followOnMessage) > 0 {
		// if the renderSystemMessage call (above) was cancelled we will
		// have a ctx.Err so make make a tidy exit
		if ctx.Err() != nil {
			return
		}

		// all good so complete the system message
		time.Sleep(1000 * time.Millisecond)
		wg.Add(1)
		if err = renderSystemMessage(ctx, conn, followOnMessage, sessionId, settings, options); err != nil {
			renderResults = append(renderResults, err)
		}
	}

	// have some errors so send them back to listener.go
	if len(renderResults) > 0 {
		chResult <- renderResults
	}
}

func renderSystemMessage(ctx context.Context, conn net.Conn, message string, sessionId string, settings config.Config, options RenderOptions) error {
	// FIXME With merged pages that are 960 char long, rendering this causes a scroll, this shouldn't happen if HOME/VTAB is used, should it?
	var (
		err         error
		cursorChars strings.Builder
		lastChar    string
	)

	logger.LogInfo.Printf("System Message: %s\r\n", message)

	// position the cursor to the bottom row, column 0
	if err = PositionCursor(conn, 0, globals.ROWS-1, !settings.Server.DisableVerticalRollOver); err != nil {
		return err
	}

	// convert the message markup and swap the pound signs
	if message, err = convert.MarkupToRawV(message); err != nil {
		return err
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
	if err = renderBuffer(ctx, conn, []byte(cursorChars.String()), settings, options); err != nil {
		return err
	}
	return nil
}

func renderHeader(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) error {

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
		if err = renderBuffer(ctx, conn, buffer, settings, options); err != nil {
			return err
		}
		return nil
	}

	//actual header text
	if len(frame.HeaderText) > 0 {
		if headerText, err = convert.MarkupToRawV(frame.HeaderText); err != nil {
			return err
		}
	} else {
		if headerText, err = convert.MarkupToRawV(settings.Server.Strings.DefaultHeaderText); err != nil {
			return err
		}
	}

	// the extra HOME characters are sent as they provide a delay for slower machines to complete a clear screen
	// it is safer to send repeated HOME than a NULL
	if !frame.DisableClear && frame.FrameType == globals.FRAME_TYPE_INITIAL {
		metadata = getMetaData(time.Now())
		cls = string(globals.CLS) + string(globals.HOME) + string(globals.HOME) + string(globals.HOME) + metadata + string(globals.CLS) + string(globals.HOME) + string(globals.HOME) + string(globals.HOME)
	} else if !frame.DisableClear {
		cls = string(globals.CLS) + string(globals.HOME) + string(globals.HOME) + string(globals.HOME)
	} else {
		cls = string(globals.HOME) + string(globals.HOME) + string(globals.HOME) + string(globals.HOME)
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

	if err = renderBuffer(ctx, conn, buffer, settings, options); err != nil {
		return err
	}
	return nil
}

func renderTitle(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) error {

	var (
		data string
		err  error
	)

	// no titles for test pages
	if frame.FrameType == globals.FRAME_TYPE_TEST {
		return nil
	}

	// split the type as this allows for comma separated params.
	dataType, rows := utils.ParseDataType(frame.Title.Type)

	// Render the content
	switch dataType {
	case "markup":
		if data, err = convert.MarkupToRawV(frame.Title.Data); err != nil {
			return err
		}
		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	case globals.CONTENT_TYPE_RAW, globals.CONTENT_TYPE_RAWV:

		data = frame.Title.Data

		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}

	case globals.CONTENT_TYPE_RAWT:
		if data, err = convert.RawTToRawV(frame.Title.Data, 0, 23, 0, 39, true); err != nil {
			return err
		}
		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}

	case globals.CONTENT_TYPE_EDITTF, "edittf", globals.CONTENT_TYPE_ZXNET:
		// for an edit.tf title field a number can be added to the type e.g. edit.tf,7 which will take the first 7
		// rows rather than the default (usually 4). This can be used to have odd frames with a larger number of rows
		// used for the title.
		if rows == 0 {
			rows = settings.Server.EditTfTitleRows
		}
		// edit.tf mode for title only returns rows 1 to 4 inclusive from the edit.tf frame.
		if data, err = convert.EdittfToRawV(frame.Title.Data, 1, rows, !settings.Server.Antiope); err != nil {
			return err
		}
		// apply any merge-data
		if frame.Title.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Title.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	}

	return nil
}

func renderContent(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) error {

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
			return err
		}
		// swap placeholders
		data = populatePlaceholders(data, settings, sessionId, options)
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	case globals.CONTENT_TYPE_RAW, globals.CONTENT_TYPE_RAWV:
		data = populatePlaceholders(frame.Content.Data, settings, sessionId, options)
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	case "rawT":
		if data, err = convert.RawTToRawV(frame.Content.Data, 0, 23, 0, 39, true); err != nil {
			return err
		}
		data = populatePlaceholders(data, settings, sessionId, options)
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	case globals.CONTENT_TYPE_EDITTF, "edittf", globals.CONTENT_TYPE_ZXNET:
		if frame.FrameType == globals.FRAME_TYPE_TEST {
			rowEnd = 23
		} else {
			rowEnd = 22
		}
		// edit.tf is teletext, so get full teletext page
		if data, err = convert.EdittfToRawV(frame.Content.Data, 1, rowEnd, !settings.Server.Antiope); err != nil {
			return err
		}
		data = populatePlaceholders(data, settings, sessionId, options)
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	}
	return nil
}

func renderFooter(ctx context.Context, conn net.Conn, frame *types.Frame, sessionId string, settings config.Config, options RenderOptions) error {

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
			return err
		}
		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}

	case globals.CONTENT_TYPE_RAW, globals.CONTENT_TYPE_RAWV:

		data = frame.Footer.Data

		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}

	case "rawT":
		if data, err = convert.RawTToRawV(frame.Content.Data, 0, 23, 0, 39, true); err != nil {
			return err
		}
		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}

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
			return err
		}
		// apply any merge-data
		if frame.Footer.MergeData != nil {
			if data, err = convert.RawVMerge(data, frame.Footer.MergeData, rows); err != nil {
				return err
			}
		}
		if err = renderBuffer(ctx, conn, []byte(data), settings, options); err != nil {
			return err
		}
	}

	return nil
}

func renderBuffer(ctx context.Context, conn net.Conn, buffer []byte, settings config.Config, options RenderOptions) error {

	var (
		err error
	)

	if settings.Server.Antiope {
		if buffer, err = convert.RawVToAntiope(buffer); err != nil {
			return err
		}
	}

	for _, b := range buffer {

		// process any requested cancellation
		select {
		case <-ctx.Done():
			// channel has a true, so cancel
			logger.LogInfo.Print("Rendering cancelled.")
			return nil // channel closed so cancel
		default:
		}

		// set EVEN parity if settings.General.Parity == true
		if settings.General.Parity {
			b = utils.SetEvenParity(b)
		}

		if _, err = conn.Write([]byte{b}); err != nil {
			logger.LogError.Print(err)
			return &NetworkError{}
		}

		// slow down to match the baud rate
		time.Sleep(time.Duration(options.BaudRate))

	}

	return nil
}

func PositionCursor(conn net.Conn, x int, y int, useRollover bool) error {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "PositionCursor")
	}

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
