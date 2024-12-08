package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/johnnewcombe/telstar-library/convert"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-library/utils"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/dal"
	"github.com/johnnewcombe/telstar/netClient"
	"github.com/johnnewcombe/telstar/renderer"
	"github.com/johnnewcombe/telstar/response"
	"github.com/johnnewcombe/telstar/routing"
	"github.com/johnnewcombe/telstar/session"
	"github.com/johnnewcombe/telstar/synchronisation"
	"net"
	"strings"
	"time"
	//"sync"
)

var listenerWg = synchronisation.WaitGroupWithCount{}

func Start(port int, settings config.Config) error {

	var (
		connectionNumber int
		conn             net.Conn
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Server.Host, port))
	if err != nil {
		return err
	}

	for {
		// blocks until an incoming connection is made
		// when a connection is made it returns a net.Conn object
		conn, err = listener.Accept()
		if err != nil {
			logger.LogError.Print(err)
			continue
		}
		connectionNumber++
		logger.LogInfo.Printf("Incoming connection [%d]!", connectionNumber)

		// handles one connection at a time
		listenerWg.Add(1)
		go handleConn(conn, connectionNumber, settings)
		logger.LogInfo.Printf("Current Session count: %3d.", listenerWg.GetCount())

	}

}

// handleConn() handles one connection at a time
func handleConn(conn net.Conn, connectionNumber int, settings config.Config) {

	var (
		reader           *bufio.Reader
		ok               bool
		telnetParser     TelnetParser
		minitelParser    MinitelParser
		routingRequest   routing.RouterRequest
		routingResponse  routing.RouterResponse
		telnetResponse   string
		minitelResponse  string
		initBytes        []byte
		currentFrame     types.Frame
		err              error
		frame            types.Frame
		hasFollowOnFrame bool
		baudRate         int
		inputByte        byte
		//deviceControlReceived bool
		currentUser       types.User
		responseFrameData response.ResponseData
		sessionId         string
		lastCharReceived  int64
		now               int64
		carouselDelay     int
		autoRefreshDelay  int
		autoRefreshFrame  bool
		networkError      *renderer.NetworkError
		logPreAmble       string
	)

	wg := synchronisation.WaitGroupWithCount{}

	// this is used to enable the rendering goroutine to be cancelled.
	ctx, cancel := context.WithCancel(context.Background())
	chResult := make(chan renderer.RenderResults)

	logPreAmble = utils.FormatLogPreAmble(session.GetSessionCount(), connectionNumber, utils.GetIpAddress(conn))

	// anonymous function used to ensure order is correct
	defer func() {

		// remove any data in the rendering result channel otherwise the renderer
		// could get blocked meaning it would never see the cancel message this
		// would mean that the following wg.Wait() call would never complete, the
		// client connection would not get closed and a session would stay alive.
		// TODO: could we use a buffered channel for results to avoid the need for this?
		select {
		case results := <-chResult:
			// channel has some data
			for _, err = range results {
				logger.LogError.Printf("%s %v", logPreAmble, err)
			}
		default:
		}

		// wait for go routines to complete
		wg.Wait() //

		// close channels
		close(chResult)

		// close the connection
		logger.LogInfo.Printf("%sClosing connection.", logPreAmble)
		if err = conn.Close(); err != nil {
			logger.LogError.Print(err)
		}

		session.DeleteSession(sessionId)
		logPreAmble = utils.FormatLogPreAmble(session.GetSessionCount(), connectionNumber, utils.GetIpAddress(conn))
		logger.LogInfo.Printf("%sSession deleted.", logPreAmble)

		//indicate to the listener that we are done
		listenerWg.Done()
	}()

	// get the corresponding user from the database,
	// FIXME this causes error on DEV Database NOTE that this is about BASE PAGE not USER ID
	//  ERROR:   2023/01/30 21:39:44 listener.go:90: finding user 7777777777: error decoding key base-page: cannot decode string into an integer type
	if currentUser, err = dal.GetUser(settings.Database.Connection, globals.GUEST_USER); err != nil {
		logger.LogError.Printf("%s: %v, has the user been created?", logPreAmble, err)
	}

	// create the users session
	sessionId = utils.CreateGuid()
	currentSession := session.CreateSession(sessionId, currentUser, connectionNumber, utils.GetIpAddress(conn))

	// update the preamble
	logPreAmble = utils.FormatLogPreAmble(session.GetSessionCount(), connectionNumber, utils.GetIpAddress(conn))
	logger.LogInfo.Printf("%sSession created with SessionId %s.", logPreAmble, sessionId)
	logger.LogInfo.Printf("%sLogged in as user: %s (%s).", logPreAmble, currentUser.Name, currentUser.UserId)

	// create a new buffered reader
	reader = bufio.NewReader(conn)

	// wait for a second or three to allow the line to settle before continuing
	// this also gives manual dial connections time to switch the modem online.
	time.Sleep(time.Second * globals.CONNECT_DELAY_SECS)

	// Removed in favour of DC parsing...
	//   see inputByte, minitelResponse = minitelParser.ParseMinitelDc(inputByte) below
	// send the <PRO1>/7B (where <PRO1> is 1B/39) to the client
	//if _, err = conn.Write([]byte(globals.MINITEL_ENQ_ROM)); err != nil {
	// logger.LogError.Printf(%d:%s: err)
	//	return
	//}

	// loop as each character is received
	for {
		//conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		if err = conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
			cancel()
			return // via defer() function
		}

		start := time.Now()
		ok, inputByte = readByte(reader)

		// process any errors from the renderer
		select {
		case results := <-chResult:
			// channel has some data
			for _, err = range results {

				logger.LogError.Printf("%s%v", logPreAmble, err)

				// cancel for network error
				if errors.As(err, &networkError) {
					// need to wait a moment as gateway connections will be shutting down.
					//time.Sleep(time.Millisecond * 500)
					cancel()
					return // via defer() function
				}
			}
		default:
		}

		if !ok {

			// each !ok read should be 500ms since the last, anything shorter means that
			// the connection has been broken
			if time.Since(start).Milliseconds() < 50 {

				//error, the connection has clearly been broken
				logger.LogInfo.Printf("%sClient has disconnected.", logPreAmble)
				cancel()
				return // TODO: should we go around again to read any results that there might have been see defer()

			}

			if currentFrame.Carousel {

				// if we are rendering or whatever e.g. wait group is > 0
				// just keep resetting the counter, otherwise increment it
				if wg.GetCount() > 0 {
					carouselDelay = 0
				} else {
					carouselDelay++
				}

				if carouselDelay == settings.Server.CarouselDelay*2 {
					// display next page
					carouselDelay = 0 // reset
					inputByte = globals.HASH
				}
			}

			if currentFrame.IsValid() && inputByte != globals.HASH {

				// we have timed out waiting for a char but we have a
				// current frame then check for auto refreshh and carousel
				// or continue look for further input

				// removed 24/11/2024
				//time.Sleep(100 * time.Millisecond)

				// if we are rendering or whatever e.g. wait group is > 0
				// just keep resetting the counter, otherwise increment it
				if wg.GetCount() > 0 {
					autoRefreshDelay = 0
				} else {
					autoRefreshDelay++
				}

				// if we have reached auto refresh time, set the flag
				if autoRefreshDelay >= settings.Server.AutoRefreshDelay*2 {
					autoRefreshFrame = true
				} else {
					// back to waiting for input
					continue
				}

			} else {
				// If the current frame is not set then this is an initial connection so
				// carry on by setting the input byte to 0x5f
				inputByte = globals.HASH
			}
		}

		if inputByte == 0 {
			continue
		}

		if !settings.Server.DisableMinitelParser {
			// pass through the Minitel parser, this will absorb any negotiation and
			// set minitelParser.MinitelState

			logger.LogInfo.Printf("%sChecking for Minitel Terminal.", logPreAmble)
			inputByte, minitelResponse = minitelParser.ParseMinitelDc(inputByte)

			// Minitel parser may need to send a response to the client, this is done here
			if len(minitelResponse) > 0 {
				if _, err = conn.Write([]byte(minitelResponse)); err != nil {
					logger.LogError.Printf("%s%v", logPreAmble, err)
				}
			}

			if minitelParser.MinitelState == MINITEL_connected && !settings.Server.Antiope {
				logger.LogInfo.Printf("%sMinitel terminal, configuring Antiope support.", logPreAmble)
				settings.Server.Antiope = true

				// If we subsequently make a connection to another service (gateway). we need to relay the contents of
				//  minitelParser.Buffer to that service.
				initBytes = minitelParser.Buffer
			}
		}

		// is this an auto refresh i.e. no inputByte
		if autoRefreshFrame {
			logger.LogInfo.Printf("%sAutomatic Refresh of frame: %d%s.", logPreAmble, currentFrame.PID.PageNumber, currentFrame.PID.FrameId)
		} else {
			logger.LogInfo.Printf("%sCharacter received: [%0x] %s (%d).", logPreAmble, inputByte, BtoA(inputByte), inputByte)
		}

		now = time.Now().UnixMilli() // nano seconds

		// TODO What does this actually do, are we talking about the user using a Viewdata client
		//  with a TNC? how would that work with CR being needed and the TNC in Line mode
		// if a hash is received within 100ms of the previous character and we are in immediate mode
		// ignore it. This should help packet radio operation where someone presses a menu choice e.g. 2
		// and hits Return (hash) as the TNC will send 2# in this scenario the hash needs to be ignored.
		if inputByte == globals.HASH && routingResponse.ImmediateMode {

			if lastCharReceived >= now-globals.HASH_GUARD_TIME {
				logger.LogInfo.Printf("%sCharacter received: [%0x] %s (%d) was received within %dms of previous character whilst in immediate mode.", logPreAmble, inputByte, BtoA(inputByte), inputByte, globals.HASH_GUARD_TIME)
				continue
			}
		}

		lastCharReceived = time.Now().UnixMilli()

		// pass through the telnet parser, this will absorb any negotiation and
		// set telnetParser.TelnetConnection to true if a telnet negotiation was detected
		// this will also absorb input bytes OD and 0A changing them to 0.
		inputByte, telnetResponse = telnetParser.ParseTelnet(inputByte, currentSession)

		// telnet parser may need to send a response to the client, this is done here
		if len(telnetResponse) > 0 {
			if _, err = conn.Write([]byte(telnetResponse)); err != nil {
				logger.LogError.Printf("%s%v", logPreAmble, err)
			}
		}

		// set the baud rate depending upon the connection type
		if telnetParser.TelnetConnection {
			logger.LogInfo.Printf("%sTelnet client, Baud rate is maximum value.", logPreAmble)
			baudRate = globals.BAUD_MAX
		} else {
			logger.LogInfo.Printf("%sBaud rate is %s.", logPreAmble, globals.BAUD_DISPLAY)
			// slow down the baud rate if not a telnet connection.
			baudRate = globals.BAUD_RATE
		}

		if settings.Server.Antiope {
			inputByte = convert.AntiopeInputTranslation(inputByte)
		}

		// ignore zero (NULL) as this is what the telnet parser returns if the character is part of a Telnet negotiation
		// however, if this is an auto refresh request inputByte will be 0 so can be ignored
		if inputByte != 0 || autoRefreshFrame {

			if settings.General.Parity {
				// added for 7E1 connections over TCP e.g. WiFi Modems connected to old 7E1 machines
				inputByte = inputByte & 0x7f
			}

			// cancel any rendering by calling the cancel method of the context
			if inputByte >= 0x30 && inputByte <= 0x39 || inputByte == globals.HASH || inputByte == globals.ASTERISK {
				cancel()
			}

			if currentFrame.FrameType == "gateway" {

				logger.LogInfo.Printf("%sCurrent frame is type Gateway.", logPreAmble)

				// to get here we must have navigated to a gateway page i.e. a page with the
				// frame type of gateway from this we can get the connection details
				// validate the address etc.
				if !currentFrame.Connection.IsValid() {
					logger.LogError.Printf("%sGateway connection details are invalid.", logPreAmble)
				}

				// send a CLS back to the client as we can't guarantee that a service will
				if _, err = conn.Write([]byte{globals.CLS}); err != nil {
					logger.LogError.Printf("%s%v", logPreAmble, err)
				}

				// returns a bool
				_ = netClient.Connect(conn, currentFrame.Connection.GetUrl(), connectionNumber, baudRate, initBytes, settings)

				// Gateway complete, so back to main index page
				routing.ForceRoute(settings.Server.Pages.GatewayErrorPage, "a", &routingResponse)

			} else if currentFrame.FrameType == "response" {

				// at this point we have rendered the frame, the currentFrame is set to the response frame
				// and we are waiting for data entry characters to complete the repsponse frame field(s) etc.
				logger.LogInfo.Printf("%sCurrent frame is type Response, passing input byte to the Response Processor.", logPreAmble)

				if err = response.Process(sessionId, conn, inputByte, &currentFrame, &responseFrameData, baudRate, settings); err != nil {

					// error with plugin
					logger.LogError.Printf("%sThe Response Processsor failed with input byte: %02x. %v", logPreAmble, inputByte, err)

					// force route to external exception error page 9902.
					routing.ForceRoute(settings.Server.Pages.ResponseErrorPage, "a", &routingResponse)

				} else if responseFrameData.Complete {
					responseFrameData.Clear()
					pid := currentFrame.ResponseData.Action.PostActionFrame
					routing.ForceRoute(pid.PageNumber, pid.FrameId, &routingResponse)
				} else {
					// not complete so continue to capture the next input byte
					continue
				}

			} else if !utils.IsValidPageId(currentFrame.GetPageId()) {

				// this would only typically be the case when the handler first starts and
				// there is no current frame
				logger.LogInfo.Printf("%sNo current frame set, using start page.", logPreAmble)

				// this is the start page so insert this into the routing process
				if settings.Server.Authentication.Required {
					routing.ForceRoute(settings.Server.Pages.LoginPage, "a", &routingResponse)
				} else {
					routing.ForceRoute(settings.Server.Pages.StartPage, "a", &routingResponse)
				}
			} else if autoRefreshFrame {
				// auto refresh of frame
				autoRefreshFrame = false
				routing.ForceRoute(currentFrame.PID.PageNumber, currentFrame.PID.FrameId, &routingResponse)

			} else {

				// current page exists so normal routing is required
				logger.LogInfo.Printf("%sCurrent frame is set (%d%s), invoking the routing process with char [%x] %s (%d).",
					logPreAmble, currentFrame.PID.PageNumber, currentFrame.PID.FrameId, inputByte, BtoA(inputByte), inputByte)

				// update the routingRequest with the byte received
				routingRequest.InputByte = inputByte

				// populate the routing request
				routingRequest.CurrentPageId = currentFrame.GetPageId()
				routingRequest.RoutingTable = currentFrame.RoutingTable
				routingRequest.HasFollowOnFrame = hasFollowOnFrame
				routingRequest.SessionId = sessionId

				// process routing
				if err = routing.ProcessRouting(&routingRequest, &routingResponse, currentSession); err != nil {
					logger.LogError.Printf("%s%v", logPreAmble, err)
				}
			}

			/* at this point the response.status is one of the following:

			UNDEFINED               = iota // 0 not actually used
			ROUTING_MESSAGE_UPDATED        // 1 - do nothing and wait for the next char
			VALID_PAGE_REQUEST             // 2 - get the page if it exists or rended PNF nav message
			INVALID_PAGE_REQUEST           // 3 - rended PNF nav message
			INVALID_CHARACTER              // 4 - log error and wait for the next char
			*/

			switch routingResponse.Status {
			case routing.Undefined:
				// 0 log error and wait for the next char
				logger.LogError.Printf("%sroutingResponse.status was set to UNDEFINED, this is an unexpected error", logPreAmble)

			case routing.RouteMessageUpdated:
				//1 - echo char and wait for the next char
				logger.LogInfo.Printf("%sRouting buffer updated with char [%x] %s (%d).", logPreAmble, inputByte, BtoA(inputByte), inputByte)

				// echo the char if the char is valid AND the waitgroup count is zero
				// i.e. nothing is being rendered
				if wg.GetCount() > 0 {
					// something is being rendered so put cursor at the bottom row
					// this bypasses the
					if err = renderer.PositionCursor(conn, 0, globals.ROWS-1, !settings.Server.DisableVerticalRollOver); err != nil {
						logger.LogError.Printf("%s%v", logPreAmble, err)
					}
				}
				if inputByte > 0x20 && inputByte < 0x7f {
					// echo the command
					// FIXME How is the cursor handled (if there is one?
					//  perhaps it should be left to the client!
					if currentFrame.FrameType != globals.FRAME_TYPE_TEST {
						if _, err = conn.Write([]byte{inputByte}); err != nil {
							logger.LogError.Printf("%s%v", logPreAmble, err)
						}
					}
				}

			case routing.BufferCharacterDeleted:

				// send BS, SPC and BS
				if _, err = conn.Write([]byte{0x08, 0x20, 0x08}); err != nil {
					logger.LogError.Printf("%s%v", logPreAmble, err)
				}

			case routing.ValidPageRequest:

				// 2 - get the page if it exists or rended PNF nav message
				// get the page

				// get the frame from db or cache

				if frame, err = getFrame(sessionId, routingResponse.NewPageId, settings); err != nil ||
					!utils.IsValidPageId(frame.GetPageId()) {

					// this means that the frame cannot be found
					// e.g. it does not exist, has visibility set to false or some other db error
					logger.LogWarn.Printf("%s%v", logPreAmble, err)

					// clear the buffer
					routingResponse.RoutingBuffer = ""

					if currentFrame.FrameType != globals.FRAME_TYPE_TEST {
						// Render PNF nav message
						var renderOptions = renderer.RenderOptions{
							BaudRate: baudRate,
						}

						ctx, cancel = context.WithCancel(context.Background())
						wg.Add(1)
						go renderer.RenderTransientSystemMessage(ctx, conn, &wg, currentFrame.NavMessageNotFound, currentFrame.NavMessage, currentSession, settings, renderOptions, chResult)
					}

				} else {
					// page id is valid, meaning that the frame was retrieved
					logger.LogInfo.Printf("%sFrame retrieved: %d%s.", logPreAmble, frame.PID.PageNumber, frame.PID.FrameId)

					// We need to determine if we have a Follow on frame
					if hasFollowOnFrame, err = existsFollowOnFrame(sessionId, frame.GetPageId(), settings); err != nil {
						logger.LogError.Print(err) // info as non-exist errors are expected
					}

					if hasFollowOnFrame {
						logger.LogInfo.Printf("%sA follow-on frame for frame %s, exists.", logPreAmble, frame.GetPageId())
					} else {
						logger.LogInfo.Printf("%sThe follow-on frame for frame %s, does not exist.", logPreAmble, frame.GetPageId())
					}

					if currentFrame.Carousel {
						logger.LogInfo.Printf("%sFrame %s is a Carousel frame.", logPreAmble, frame.GetPageId())
					}

					if !routingResponse.HistoryPage {
						logger.LogInfo.Printf("%sAdding frame %s to history.", logPreAmble, frame.GetPageId())
						session.PushHistory(sessionId, frame.GetPageId())
					}

					routingResponse.Clear()

					logger.LogInfo.Printf("%sRendering frame %s.", logPreAmble, frame.GetPageId())

					// Render the frame
					var renderOptions = renderer.RenderOptions{
						HasFollowOnFrame: hasFollowOnFrame,
						//ClearScreen:      true,
						BaudRate: baudRate,
					}

					// make sure we have some nav messages.
					if len(frame.NavMessageNotFound) == 0 {
						frame.NavMessageNotFound = settings.Server.Strings.DefaultPageNotFoundMessage
					}
					if len(frame.NavMessage) == 0 {
						frame.NavMessage = settings.Server.Strings.DefaultNavMessage
					}

					// create a new context that can be used to allow rendering to be cancelled
					ctx, cancel = context.WithCancel(context.Background())
					wg.Add(1)
					go renderer.Render(ctx, conn, &wg, &frame, currentSession, settings, renderOptions, chResult)

					if frame.FrameType == "exit" {
						cancel()
						return
					}
					currentFrame = frame

				}

			case routing.InvalidPageRequest:

				if currentFrame.FrameType != globals.FRAME_TYPE_TEST {

					var renderOptions = renderer.RenderOptions{
						BaudRate: baudRate,
					}

					ctx, cancel = context.WithCancel(context.Background())
					wg.Add(1)
					go renderer.RenderTransientSystemMessage(ctx, conn, &wg, currentFrame.NavMessageNotFound, currentFrame.NavMessage, currentSession, settings, renderOptions, chResult)
				}

			case routing.InvalidCharacter:
				// 4 - log warning and wait for the next char
				logger.LogWarn.Printf("%sAn invalid character was received from the connected client [%0x] %s (%d).", logPreAmble, inputByte, BtoA(inputByte), inputByte)
			}
		}
	}
}

/*
func readByteNew(conn net.Conn) (bool, byte) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "readByteNew")
	}

	// FIXME: Why do we have this and it not be used? This uses Read not ReadByte
	//   and includes a timeout
	var (
		inputByte byte
		err       error
	)

	buf := make([]byte, 1)
	if err = conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
		return false, 0
	}
	if _, err = conn.Read(buf); err != nil {
		return false, 0
	}

	// we do not need buffered input as we are only expecting single char inputs
	inputByte = buf[0]
	if inputByte > 0 {
		logger.LogInfo.Printf("@")
	} else {
		return false, 0
	}
	return true, inputByte
}
*/

func readByte(reader *bufio.Reader) (bool, byte) {

	// get a byte
	inputByte, err := reader.ReadByte()
	if err != nil {
		return false, inputByte
	}
	return true, inputByte
}

func existsFollowOnFrame(sessionId, pageId string, settings config.Config) (bool, error) {

	var (
		followOnFrame   types.Frame
		err             error
		followOnFrameId string
	)

	if followOnFrameId, err = routing.GetFollowOnPageId(pageId); err != nil {
		return false, err
	} else {
		// frame has a follow on frame id so tray and get that frame
		if followOnFrame, err = getFrame(sessionId, followOnFrameId, settings); err != nil {
			return false, nil
		} else if !utils.IsValidPageId(followOnFrame.GetPageId()) {
			return false, nil
		} else {
			return true, nil
		}
	}
}

func getFrame(sessionId string, pageId string, settings config.Config) (types.Frame, error) {

	var (
		pageNumber int
		frameId    string
		err        error
		frame      types.Frame
	)

	if pageNumber, frameId, err = utils.ConvertPageIdToPID(pageId); err != nil {
		return frame, err
	}
	primary := strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	if frame, err = session.GetFrameFromCache(sessionId, pageId); err != nil ||
		!utils.IsValidPageId(frame.GetPageId()) {
		if frame, err = dal.GetFrame(settings.Database.Connection, pageNumber, frameId, primary, true); err != nil {
			return frame, err
		}
	}

	if len(frame.Redirect.FrameId) > 0 {
		// must be a redirect so re-enter this function
		if frame, err = getFrame(sessionId, frame.GetRedirectPageId(), settings); err != nil {
			return frame, err
		}
	}

	return frame, nil
}

func BtoA(b byte) string {
	// used for display purposes only
	if b >= 20 && b < 128 {
		return string(b)
	} else {
		return "."
	}
}
