package server

import (
	"bufio"
	"context"
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

func Start(port int, settings config.Config) error {

	logger.LogInfo.Printf("Starting Videotex Server %s on port %d", settings.Server.DisplayName, port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Server.Host, port))

	if err != nil {
		return err
	}

	for {
		// blocks until an incomming connection is made
		// when a connectioj is made it returns a net.Conn object
		conn, err := listener.Accept()
		if err != nil {
			logger.LogError.Print(err)
			continue
		}
		logger.LogInfo.Print("Incoming connection!")

		// handles one connection at a time
		go handleConn(conn, settings)

	}
}

// handleConn() handles one connection at a time
func handleConn(conn net.Conn, settings config.Config) {

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
	)

	defer closeConn(conn)

	// this is used to enable the rendering goroutine to be cancelled.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//wg = sync.WaitGroup{}
	wg := synchronisation.WaitGroupWithCount{}
	defer wg.Wait()

	// get the corresponding user from the database,
	// FIXME this causes error on DEV Database NOTE that this is about BASE PAGE not USER ID
	//  ERROR:   2023/01/30 21:39:44 listener.go:90: finding user 7777777777: error decoding key base-page: cannot decode string into an integer type
	if currentUser, err = dal.GetUser(settings.Database.Connection, globals.GUEST_USER); err != nil {
		logger.LogError.Print(err)
	}
	logger.LogInfo.Printf("Logged in as user: %s (%s)", currentUser.Name, currentUser.UserId)

	// create the users session
	sessionId = utils.CreateGuid()
	session.CreateSession(sessionId, currentUser)
	defer session.DeleteSession(sessionId)

	// TODO Ensure that the session Id is displayed to the user somehow, perhaps on the session page?
	//  all logging should include the session Id or a hash/CRC or something e.g.
	//  sess id: 33efccb148144be49695e407239c4124
	//  it doesn't need to be perfect so could use 33ef-4124 i.e. first 4 and last 4 chars ??

	logger.LogInfo.Printf("Session created with SessionId %s\r\n", sessionId)

	// a simple flag to indicate that this is the first time through the loop
	//initialPass = true
	//deviceControlReceived = false

	// create a new buffered reader
	reader = bufio.NewReader(conn)

	// wait for a second or three to allow the line to settle before continuing
	// this also gives manual dial connections time to switch the modem online.
	time.Sleep(time.Second * globals.CONNECT_DELAY_SECS)

	// loop as each character is received
	for {

		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		ok, inputByte = readByte(reader)

		//} else
		if !ok {

			// if we timeout waiting for a char but we have a current frame then
			// wait 100ms and look for further input.
			if len(currentFrame.PID.FrameId) > 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				// If the current frame is not set then this is an initial connection so
				// carry on by setting the input byte to 0x5f
				inputByte = globals.HASH
			}
		}

		logger.LogInfo.Printf("Character received: [%0x] %s (%d)", inputByte, BtoA(inputByte), inputByte)

		now = time.Now().UnixMilli() // nano seconds

		// TODO What does this actually do, are we talking about the user using a Viewdata client
		//  with a TNC? how would that work with CR being needed and the TNC in Line mode
		// if a hash is received within 100ms of the previous character and we are in immediate mode
		// ignore it. This should help packet radio operation where someone presses a menu choice e.g. 2
		// and hits Return (hash) as the TNC will send 2# in this scenario the hash needs to be ignored.
		if inputByte == globals.HASH && routingResponse.ImmediateMode {

			if lastCharReceived >= now-globals.HASH_GUARD_TIME {
				logger.LogInfo.Printf("Character received: [%0x] %s (%d) was received within %dms of previous character whilst in immediate mode.", inputByte, BtoA(inputByte), inputByte, globals.HASH_GUARD_TIME)
				continue
			}
		}

		lastCharReceived = time.Now().UnixMilli()

		// pass through the Minitel parser, this will absorb any negotiation and
		// set minitelParser.MinitelConnection to true if a Minitel negotiation was detected
		inputByte, minitelResponse = minitelParser.ParseMinitel(inputByte)

		// Minitel parser may need to send a response to the client, this is done here
		if len(minitelResponse) > 0 {
			if _, err = conn.Write([]byte(minitelResponse)); err != nil {
				logger.LogError.Print(err)
			}
		}

		if minitelParser.MinitelConnection {
			logger.LogInfo.Print("Minitel terminal, configuring Antiope support.")
			settings.Server.Antiope = true

			// If we subsequently make a connection to another service (gateway). we need to relay the contents of
			//  minitelParser.Buffer to that service.
			initBytes = minitelParser.Buffer
		}

		// pass through the telnet parser, this will absorb any negotiation and
		// set telnetParser.TelnetConnection to true if a telnet negotiation was detected
		// this will also absorb input bytes OD and 0A changing them to 0.
		inputByte, telnetResponse = telnetParser.ParseTelnet(inputByte)

		// telnet parser may need to send a response to the client, this is done here
		if len(telnetResponse) > 0 {
			if _, err = conn.Write([]byte(telnetResponse)); err != nil {
				logger.LogError.Print(err)
			}
		}

		// set the baud rate depending upon the connection type
		if telnetParser.TelnetConnection {
			logger.LogInfo.Print("Telnet client, Baud rate is maximum value.")
			baudRate = globals.BAUD_MAX
		} else {
			logger.LogInfo.Print("Baud rate is 2400 baud.")
			// slow down the baud rate if not a telnet connection.
			baudRate = globals.BAUD_RATE
		}

		if settings.Server.Antiope {
			inputByte = convert.AntiopeInputTranslation(inputByte)
		}

		// ignore zero (NULL) as this is what the telnet parser returns if the character is part of a Telnet negotiation
		if inputByte != 0 {

			if settings.General.Parity {
				// remove parity
				inputByte = inputByte & 0x7f
			}

			// cancel any rendering by calling the cancel method of the context
			if inputByte >= 0x30 && inputByte <= 0x39 || inputByte == globals.HASH || inputByte == globals.ASTERISK {
				cancel()
			}

			if currentFrame.FrameType == "gateway" {

				logger.LogInfo.Print("Current frame is type Gateway.")

				// to get here we must have navigated to a gateway page i.e. a page with the
				// frame type of gateway from this we can get the connection details
				// validate the address etc.
				if currentFrame.Connection.IsValid() {
					logger.LogError.Print("Gateway connection details are invalid.")
				}

				// send a CLS back to the client as we can't guarantee that a service will
				if _, err := conn.Write([]byte{globals.CLS}); err != nil {
					logger.LogError.Println(err)
				}

				netClient.Connect(conn, currentFrame.Connection.GetUrl(), settings.Server.DLE, baudRate, initBytes)

				// Gateway complete, so back to main index page
				routing.ForceRoute(settings.Server.Pages.GatewayErrorPage, "a", &routingRequest, &routingResponse)

			} else if currentFrame.FrameType == "response" {

				// at this point we have rendered the frame, the currentFrame is set to the response frame
				// and we are waiting for data entry characters to complete the repsponse frame field(s) etc.
				logger.LogInfo.Print("Current frame is type Response, passing input byte to the Response Processor.")

				if err = response.Process(sessionId, conn, inputByte, &currentFrame, &responseFrameData, baudRate, settings); err != nil {

					// error with plugin
					logger.LogError.Printf("The Response Processsor failed with input byte: %02x. %v", inputByte, err)

					// force route to external exception error page 9902.
					routing.ForceRoute(settings.Server.Pages.ResponseErrorPage, "a", &routingRequest, &routingResponse)

				} else if responseFrameData.Complete {
					responseFrameData.Clear()
					pid := currentFrame.ResponseData.Action.PostActionFrame
					routing.ForceRoute(pid.PageNumber, pid.FrameId, &routingRequest, &routingResponse)
				} else {
					// not complete so continue to capture the next input byte
					continue
				}

			} else if !utils.IsValidPageId(currentFrame.GetPageId()) {

				// this would only typically be the case when the handler first starts and
				// there is no current frame
				logger.LogInfo.Print("No current frame set, using start page.")

				// this is the start page so insert this into the routing process
				if settings.Server.Authentication.Required {
					routing.ForceRoute(settings.Server.Pages.LoginPage, "a", &routingRequest, &routingResponse)
				} else {
					routing.ForceRoute(settings.Server.Pages.StartPage, "a", &routingRequest, &routingResponse)
				}
			} else {

				// current page exists so normal routing is required
				logger.LogInfo.Printf("Current frame is set (%d%s), invoking the routing process with char [%x] %s (%d).",
					currentFrame.PID.PageNumber, currentFrame.PID.FrameId, inputByte, BtoA(inputByte), inputByte)

				// update the routingRequest with the byte received
				routingRequest.InputByte = inputByte

				// populate the routing request
				routingRequest.CurrentPageId = currentFrame.GetPageId()
				routingRequest.RoutingTable = currentFrame.RoutingTable
				routingRequest.HasFollowOnFrame = hasFollowOnFrame
				routingRequest.SessionId = sessionId

				// process routing
				if err = routing.ProcessRouting(&routingRequest, &routingResponse); err != nil {
					logger.LogError.Print(err)
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
				logger.LogError.Print("routingResponse.status was set to UNDEFINED, this is an unexpected error")

			case routing.RouteMessageUpdated:
				//1 - echo char and wait for the next char
				logger.LogInfo.Printf("Routing buffer updated with char [%x] %s (%d)", inputByte, BtoA(inputByte), inputByte)

				// echo the char if the char is valid AND the waitgroup count is zero
				// i.e. nothing is being rendered
				if wg.GetCount() > 0 {
					// something is being rendered so put cursor at the bottom row
					// this bypasses the
					if err := renderer.PositionCursor(conn, 0, globals.ROWS-1, !settings.Server.DisableVerticalRollOver); err != nil {
						logger.LogError.Println(err)
					}
				}
				if inputByte > 0x20 && inputByte < 0x7f {
					// echo the command
					// FIXME How is the cursor handled (if there is one?
					//  perhaps it should be left to the client!
					if currentFrame.FrameType != globals.FRAME_TYPE_TEST {
						if _, err := conn.Write([]byte{inputByte}); err != nil {
							logger.LogError.Println(err)
						}
					}
				}

			case routing.BufferCharacterDeleted:

				// send BS, SPC and BS
				if _, err := conn.Write([]byte{0x08, 0x20, 0x08}); err != nil {
					logger.LogError.Println(err)
				}

			case routing.ValidPageRequest:

				// 2 - get the page if it exists or rended PNF nav message
				// get the page

				// get the frame from db or cache

				if frame, err = getFrame(sessionId, routingResponse.NewPageId, settings); err != nil ||
					!utils.IsValidPageId(frame.GetPageId()) {

					// this means that the frame cannot be found
					// e.g. it does not exist, has visibility set to false or some other db error
					logger.LogError.Print(err)

					// clear the buffer
					routingResponse.RoutingBuffer = ""

					if currentFrame.FrameType != globals.FRAME_TYPE_TEST {
						// Render PNF nav message
						var renderOptions = renderer.RenderOptions{
							BaudRate: baudRate,
						}

						ctx, cancel = context.WithCancel(context.Background())
						wg.Add(1)
						go renderer.RenderTransientSystemMessage(ctx, conn, &wg, currentFrame.NavMessageNotFound, currentFrame.NavMessage, settings, renderOptions)
					}

				} else {
					// page id is valid, meaning that the frame was retrieved
					logger.LogInfo.Printf("Frame retrieved: %d%s", frame.PID.PageNumber, frame.PID.FrameId)

					// We need to determine if we have a Follow on frame
					if hasFollowOnFrame, err = existsFollowOnFrame(sessionId, frame.GetPageId(), settings); err != nil {
						logger.LogInfo.Print(err) // info as non-exist errors are expected
					}

					if hasFollowOnFrame {
						logger.LogInfo.Printf("A follow-on frame for frame %s, exists.", frame.GetPageId())
					} else {
						logger.LogInfo.Printf("The follow-on frame for frame %s, does not exist.", frame.GetPageId())
					}

					logger.LogInfo.Printf("Adding frame %s to history.", frame.GetPageId())

					if !routingResponse.HistoryPage {
						session.PushHistory(sessionId, frame.GetPageId())
					}

					routingResponse.Clear()

					logger.LogInfo.Printf("Rendering frame %s.", frame.GetPageId())

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
					go renderer.Render(ctx, conn, &wg, &frame, sessionId, settings, renderOptions)

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
					go renderer.RenderTransientSystemMessage(ctx, conn, &wg, currentFrame.NavMessageNotFound, currentFrame.NavMessage, settings, renderOptions)
				}

			case routing.InvalidCharacter:
				// 4 - log warning and wait for the next char
				logger.LogWarn.Print("An invalid character was received from the connected client.")
			}
		}
	}
}

func closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		logger.LogError.Print(err)
	}
	logger.LogInfo.Print("Closing connection!")
}

func readByteNew(conn net.Conn) (bool, byte) {

	var (
		inputByte byte
	)

	buf := make([]byte, 1)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	conn.Read(buf)

	// we do not need buffered input as we are only expecting single char inputs
	inputByte = buf[0]
	if inputByte > 0 {
		logger.LogInfo.Printf("@")
	} else {
		return false, 0
	}
	return true, inputByte
}

func readByte(reader *bufio.Reader) (bool, byte) {

	// get a byte
	inputByte, err := reader.ReadByte()
	if err != nil {
		return false, inputByte
	}
	logger.LogInfo.Println("character read:", inputByte)
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
			return false, err
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
