package main

import (
	"context"
	_ "embed"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/johnnewcombe/telstar-client/comms"
	"github.com/johnnewcombe/telstar-client/constants"
	"github.com/johnnewcombe/telstar-client/display"
	"github.com/johnnewcombe/telstar-client/keyboard"
	"gopkg.in/yaml.v3"
	"image/color"
	"io/fs"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var debug bool

//go:embed connection-files/telstar.yml
var telstarConnection string

//go:embed connection-files/nxtel.yml
var nxtelConnection string

//go:embed connection-files/ccl4.yml
var ccl4Connection string

//go:embed connection-files/teefax.yml
var teefaxConnection string

//go:embed connection-files/eol-bbs.yml
var eolBbsConnection string

var configDir string

func main() {

	var (
		//demo             bool = false
		err                   error
		address               string
		fullScreen            bool
		noToolbar             bool
		startupDelay          int
		fullScreenDelay       int
		textSize              float64
		endpoint              comms.Endpoint
		lastReadFileLocation  string
		lastWriteFileLocation string
		commsClient           comms.CommunicationClient
		ctxCommsClient        context.Context
		cancelRead            context.CancelFunc
		ctxDisplay            context.Context
		cancelDisplay         context.CancelFunc
		vbox                  *fyne.Container
		scr1Container         *fyne.Container
		screen                display.Screen
		wgComms               sync.WaitGroup
		wgDisp                sync.WaitGroup
		a                     fyne.App
		w                     fyne.Window
		theme                 display.StandardTheme
		flags                 *flag.FlagSet
		port                  int
		baud                  int
	)

	flags = flag.NewFlagSet("telstar-client", flag.ExitOnError)
	flags.StringVar(&address, "address", "", "Endpoint definition file e.g. Telstar.yml.")
	flags.BoolVar(&debug, "debug", false, "Outputs debug information to standard output.")
	flags.Float64Var(&textSize, "text-size", 23.0, "Text size, can be used to counter display issues.")
	flags.BoolVar(&fullScreen, "full-screen", false, "Full screen mode (experimental).")
	flags.IntVar(&fullScreenDelay, "full-screen-delay", 5, "Full screen mode delay (experimental).")
	flags.BoolVar(&noToolbar, "no-toolbar", false, "Hides the toolbar and status line.")
	flags.IntVar(&startupDelay, "startup-delay", 0, "Delays startup of the application..")

	if err = flags.Parse(os.Args[1:]); err != nil {
		displayError(err, w)
		os.Exit(1)
	}

	if startupDelay > 0 {
		time.Sleep(time.Second * time.Duration(startupDelay))
	}

	a = app.New()

	if fullScreen {
		display.Gray = color.Gray{Y: 0}
	} else {
		display.Gray = color.Gray{Y: 32}
	}
	theme = display.StandardTheme{}

	a.Settings().SetTheme(theme)
	a.SetIcon(constants.AppIcon)

	// add the main window and resize
	w = a.NewWindow(constants.Version)
	w.CenterOnScreen()

	// pop all of the embedded connection files in the config folder
	if err = saveConnectionFilesToConfig(); err != nil {
		// if this fails we could probably just ignore it
		//displayError(err, w)
	}

	// create a wait group and make sure we wait for all goroutines to end before exiting
	wgComms = sync.WaitGroup{}
	wgDisp = sync.WaitGroup{}

	// create a  context for the screen
	ctxDisplay, cancelDisplay = context.WithCancel(context.Background())
	//ctxCommsClient, cancelRead = context.WithCancel(context.Background())

	// set up the screen
	screen = display.Screen{}
	screen.Debug = debug
	scr1Container = screen.Initialise(ctxDisplay, &wgDisp, w.Canvas(), textSize)
	status := canvas.NewText("", color.White)

	// this needs to be here incase the initial open fails, and the user selects another
	//ctxCommsClient, cancelRead = context.WithCancel(context.Background())
	// define the Open function
	openFunc := func() error {

		// create the communications client
		if endpoint.IsSerial() {
			commsClient = &comms.SerialClient{}
		} else {
			commsClient = &comms.NetClient{}
		}

		if err = commsClient.Open(endpoint); err != nil {

			return err
		}

		wgComms.Add(1)
		ctxCommsClient, cancelRead = context.WithCancel(context.Background())
		go commsClient.Read(ctxCommsClient, &wgComms, func(ok bool, b byte) {
			if ok {
				if err = screen.Write(b); err != nil {
					// It is safe to updated Fyne UI components from within a go routine, in fact its probably the only
					// way to do it and is sited in Fyne examples. See https://developer.fyne.io/started/updating
					displayError(err, w)
					return
				}
			} else {
				status.Text = "Offline"
				status.Refresh()
			}
		})
		status.Text = fmt.Sprintf("Online [%s]", endpoint.Name)
		status.Refresh()
		return nil
	}

	closeFunc := func() {
		// close the previous client and stop the read goroutine.
		// The commsClient.Read() goroutine blocks on serial/net. Closing the
		// connection/port will cause a read error and allow the go routine to continue
		// monitoring for ctx cancel
		commsClient.Close()

		// (it raises and error instead) and it is now looping until cancelled,
		// so lets cancel it
		if cancelRead != nil {
			cancelRead()
		}

		// wait for all goroutines to stop
		wgComms.Wait()
	}

	exitFunc := func() {
		// save endpoint to .recent, this will get loaded the next time we start the client.
		if err = SaveEndpoint(path.Join(configDir, comms.DEFAULT_CONNECTION), endpoint); err != nil {
			// we can ignore this maybe
			//displayError(errors.New("unable to save to the '.recent' file"), w)
		}
		w.Close()
		cancelDisplay()
		// order is important
		commsClient.Close()
		if cancelRead != nil {
			cancelRead()
		}
		wgComms.Wait()
		wgDisp.Wait()
	}

	tb := widget.NewToolbar(
		widget.NewToolbarAction(theme.Icon("cloud"),
			func() {

				ipaddressWidget := widget.NewEntry()
				ipaddressWidget.SetPlaceHolder("")
				portWidget := widget.NewEntry()
				portWidget.SetPlaceHolder("")
				telnetWidget := widget.NewCheck("", nil)
				//				widget.NewFormItem("Port", portWidget)
				initWidget := widget.NewEntry()
				initWidget.SetPlaceHolder("e.g. '5f' or '0d0a' etc.")

				input := dialog.NewForm("TCP Connection", "OK", "Cancel", []*widget.FormItem{

					widget.NewFormItem("IP Address", ipaddressWidget),
					widget.NewFormItem("TCP Port", portWidget),
					widget.NewFormItem("Requires Telnet Negotiation", telnetWidget),
					widget.NewFormItem("Initialisation Bytes", initWidget),
				}, func(confirm bool) {

					if confirm {

						name := fmt.Sprintf("%s:%s", ipaddressWidget.Text, portWidget.Text)

						endpoint = comms.Endpoint{Name: name}
						endpoint.Address.Host = ipaddressWidget.Text
						if port, err = strconv.Atoi(portWidget.Text); err != nil {
							port = 0 // this causes a 'invalid connection error' below.
						}
						endpoint.Address.Port = port
						endpoint.Init.Telnet = telnetWidget.Checked
						var hexString []byte
						if hexString, err = hex.DecodeString(initWidget.Text); err != nil {
							// ignore
						}
						endpoint.Init.InitChars = hexString

						// close the previous client and stop the read goroutine.
						closeFunc()

						// close the previous client and stop the read goroutine.
						// The commsClient.Read() goroutine blocks on serial/net. Closing the
						// connection/port will cause a read error and allow the go routine to continue
						// monitoring for ctx cancel
						//commsClient.Close()

						// (it raises and error instead) and it is now looping until cancelled,
						// so lets cancel it
						//if cancelRead != nil {
						//	cancelRead()
						//}

						// wait for all goroutines to stop
						//wgComms.Wait()

						// open the connection and run the goroutine that receives data
						if err = openFunc(); err != nil {
							status.Text = "Offline"
							status.Refresh()
							displayError(errors.New("specified connection is invalid"), w)
						} else {

							// save endpoint to .recent, this will get loaded the next time we start the client.
							if err = SaveEndpoint(path.Join(configDir, comms.DEFAULT_CONNECTION), endpoint); err == nil {

								fd := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {

									if err != nil {
										displayError(err, w)
										status.Text = "Offline"
										status.Refresh()
										return
									}

									if write != nil {

										filename := write.URI()
										path := filename.Path()

										// can't use this as a file is creatred before we get here, if we do this we end up
										// with two files, one with the extension and one without
										// if !strings.HasSuffix(path, ".yml") {
										// 	path = path + ".yml"
										// }
										if err = SaveEndpoint(path, endpoint); err != nil {
											displayError(errors.New("unable to save endpoint details"), w)
										}

									}
								}, w)

								if len(lastWriteFileLocation) == 0 {
									lastWriteFileLocation = configDir
								}

								fileURI := storage.NewFileURI(lastWriteFileLocation)
								fileLister, _ := storage.ListerForURI(fileURI)
								fd.SetLocation(fileLister)
								fd.SetFileName("untitled.yml")
								filter := storage.NewExtensionFileFilter([]string{".yml"})
								fd.SetFilter(filter)
								fd.Resize(fyne.Size{w.Canvas().Size().Width * 0.8, w.Canvas().Size().Height * 0.8})
								fd.Show()

							}
						}
					}
				}, w)

				input.Resize(fyne.Size{
					Width:  float32(math.Min(float64(w.Canvas().Size().Width*0.75), 500)),
					Height: float32(math.Min(float64(w.Canvas().Size().Height*0.3), 200)),
				})

				ipaddressWidget.Text = endpoint.Address.Host
				portWidget.Text = strconv.Itoa(endpoint.Address.Port)
				telnetWidget.Checked = endpoint.Init.Telnet
				input.Show()
			}),
		widget.NewToolbarAction(theme.Icon("serial"),
			func() {

				deviceWidget := widget.NewEntry()
				deviceWidget.SetPlaceHolder("e.g. COM1 or /dev/ttyS0 etc.")
				baudWidget := widget.NewEntry()
				baudWidget.SetPlaceHolder("e.g. 1200")
				parityWidget := widget.NewCheck("", nil)
				//				widget.NewFormItem("Port", portWidget)
				modemInitWidget := widget.NewEntry()
				modemInitWidget.SetPlaceHolder("e.g. ATDT01756664433")

				input := dialog.NewForm("Serial Connection", "OK", "Cancel", []*widget.FormItem{

					// fixme add additional widgets
					// FIXME what about serial connections

					widget.NewFormItem("Serial Device", deviceWidget),
					widget.NewFormItem("Baud Rate", baudWidget),
					widget.NewFormItem("Parity (7 bit Even)", parityWidget),
					widget.NewFormItem("Modem Initialisation String", modemInitWidget),
				}, func(confirm bool) {

					if confirm {

						name := fmt.Sprintf("%s:%s", deviceWidget.Text, baudWidget.Text)

						endpoint = comms.Endpoint{Name: name}
						endpoint.Serial.Port = deviceWidget.Text
						if baud, err = strconv.Atoi(baudWidget.Text); err != nil {
							baud = 0 // this causes a 'invalid connection error' below.
						}
						endpoint.Serial.Baud = baud
						endpoint.Serial.Parity = parityWidget.Checked
						endpoint.Serial.ModemInit = modemInitWidget.Text

						// close the previous client and stop the read goroutine.
						closeFunc()

						// close the previous client and stop the read goroutine.
						// The commsClient.Read() goroutine blocks on serial/net. Closing the
						// connection/port will cause a read error and allow the go routine to continue
						// monitoring for ctx cancel
						//commsClient.Close()

						// (it raises and error instead) and it is now looping until cancelled,
						// so lets cancel it
						//if cancelRead != nil {
						//	cancelRead()
						//}

						// wait for all goroutines to stop
						//wgComms.Wait()

						// open the connection and run the goroutine that receives data
						if err = openFunc(); err != nil {
							status.Text = "Offline"
							status.Refresh()
							displayError(errors.New("specified connection is invalid"), w)
						} else {
							// save endpoint to .recent, this will get loaded the next time we start the client.
							if err = SaveEndpoint(path.Join(configDir, comms.DEFAULT_CONNECTION), endpoint); err == nil {
								fd := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {

									if err != nil {
										displayError(err, w)
										status.Text = "Offline"
										status.Refresh()
										return
									}

									if write != nil {

										filename := write.URI()
										path := filename.Path()

										// can't use this as a file is created before we get here, if we do this we end up
										// with two files, one with the extension and one without
										// if !strings.HasSuffix(path, ".yml") {
										// 	path = path + ".yml"
										// }
										if err = SaveEndpoint(path, endpoint); err != nil {
											displayError(errors.New("unable to save endpoint details"), w)
										}

									}
								}, w)

								if len(lastWriteFileLocation) == 0 {
									lastWriteFileLocation = configDir
								}

								fileURI := storage.NewFileURI(lastWriteFileLocation)
								fileLister, _ := storage.ListerForURI(fileURI)
								fd.SetLocation(fileLister)
								fd.SetFileName("untitled.yml")
								filter := storage.NewExtensionFileFilter([]string{".yml"})
								fd.SetFilter(filter)
								fd.Resize(fyne.Size{w.Canvas().Size().Width * 0.8, w.Canvas().Size().Height * 0.8})
								fd.Show()
							}
						}
					}
				}, w)

				input.Resize(fyne.Size{
					Width:  float32(math.Min(float64(w.Canvas().Size().Width*0.75), 500)),
					Height: float32(math.Min(float64(w.Canvas().Size().Height*0.3), 200)),
				})

				deviceWidget.Text = endpoint.Serial.Port
				baudWidget.Text = strconv.Itoa(endpoint.Serial.Baud)
				parityWidget.Checked = endpoint.Serial.Parity
				input.Show()
			}),
		widget.NewToolbarAction(theme.Icon("folderOpen"),
			func() {

				fd := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {

					if err != nil {
						displayError(err, w)
						status.Text = "Offline"
						status.Refresh()
						return
					}

					if read != nil {

						fulPath, _ := filepath.Abs(read.URI().Path())
						lastReadFileLocation = filepath.Dir(fulPath)

						filename := read.URI()
						if endpoint, err = comms.GetEndpoint(filename.Path()); err != nil {
							status.Text = "Offline"
							status.Refresh()
							displayError(err, w)
							return
						}

						// close the previous client and stop the read goroutine.
						// The commsClient.Read() goroutine blocks on serial/net. Closing the
						// connection/port will cause a read error and allow the go routine to continue
						// monitoring for ctx cancel
						commsClient.Close()

						// now that the commsClient is closed the read is no longer blocking
						// (it raises and error instead) and it is now looping until cancelled,
						// so lets cancel it
						if cancelRead != nil {
							cancelRead()
						}

						// wait for all goroutines to stop
						wgComms.Wait()

						// open the connection and run the goroutine that receives data
						if err = openFunc(); err != nil {
							status.Text = "Offline"
							status.Refresh()
							displayError(errors.New("specified connection is invalid"), w)
						}

					}
				}, w)

				if len(lastReadFileLocation) == 0 {
					lastReadFileLocation = configDir
				}

				fileURI := storage.NewFileURI(lastReadFileLocation)
				fileLister, _ := storage.ListerForURI(fileURI)
				fd.SetLocation(fileLister)

				filter := storage.NewExtensionFileFilter([]string{".yml"})
				fd.SetFilter(filter)
				fd.Resize(fyne.Size{w.Canvas().Size().Width * 0.8, w.Canvas().Size().Height * 0.8})
				fd.Show()
			}),
		widget.NewToolbarAction(theme.Icon("close"),
			func() {
				closeFunc()
			}),
		widget.NewToolbarSpacer(),
		//widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.Icon("visibility"),
			func() {
				screen.Reveal()
			}),
		widget.NewToolbarAction(theme.Icon("visibilityOff"),
			func() {
				screen.Conceal()
			}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.Icon("help"),
			func() {

				richText := widget.NewRichTextFromMarkdown(constants.HelpMsg)
				richText.Wrapping = fyne.TextWrapWord
				richText.Scroll = container.ScrollVerticalOnly
				cd := dialog.NewCustom(constants.Version, "OK", richText, w)
				cd.Resize(w.Canvas().Size())
				cd.Show()
			}),

		widget.NewToolbarAction(theme.Icon("cancel"),
			func() {
				// use exitFunc rather than defer
				exitFunc()
				return
			}),
	)

	// set up the main gui canvas areas
	if fullScreen || noToolbar {
		// fixme removing the toolbar caused the thr top line to be partially hidden
		//  maybe add some kind of spacer/top padding?
		tb.Hide()
		status.Hide()
	}
	vbox = container.NewBorder(tb, status, nil, nil, scr1Container)

	// window
	//w.Canvas().Scale()
	//w.Resize(fyne.NewSize(800, 600 + tb.Size().Height + status.Size().Height))
	w.SetFixedSize(true)
	if fullScreen {
		go setFullScreen(w, fullScreenDelay)
	}
	//w.SetFullScreen(fullScreen)
	w.CenterOnScreen()
	w.SetContent(vbox)
	w.SetCloseIntercept(func() {
		exitFunc()
		return
	})

	// open the url as specified (or default if not)
	if len(address) > 0 {

		if !strings.HasSuffix(address, ".yml") {
			address = address + ".yml"
		}

		// FIXME if a full path is specified then we should not use the config folder ?
		if endpoint, err = comms.GetEndpoint(path.Join(configDir, address)); err != nil {
			status.Text = "Offline"
			status.Refresh()
			displayError(err, w)
		}

	} else {

		// looking for the endpoint in the file '.recent'
		if endpoint, err = comms.GetEndpoint(path.Join(configDir, comms.DEFAULT_CONNECTION)); err != nil {
			status.Text = "Offline"
			status.Refresh()

			// no good so far so connect to Telstar
			endpoint = comms.GetDefaultEndPoint()
		}
	}

	// open the connection and run the goroutine that receives data
	// the call to open here is executed when the program first starts
	// if this is set to start automatically on system boot, the network
	// may not be available, if it fails retry after a wait period
	for retries := 0; retries < 3; retries++ {
		if err = openFunc(); err != nil {
			time.Sleep(time.Second * 2)
			continue
		}
		break
	}

	if err != nil {
		status.Text = "Offline"
		status.Refresh()
		displayError(errors.New("specified connection is invalid"), w)
	}

	keyb := keyboard.Keyboard{}
	keyb.Initialise(w, func(b byte) {
		//fmt.Printf("Key: %h", b)

		if b == 0xB1 { // F1
			screen.Reveal()
		} else if b == 0xB2 { // F2
			screen.Conceal()
		} else if b == 0x6 {
			w.SetFullScreen(true)
		} else if b == 0x1b {
			// FIXME because we have increased the text size as part of going full screen
			//  the minimum size canot be made smaller which means that while we can come out
			//  of full screen mode we are left with a windo that we cant resize
			w.SetFullScreen(false)
		} else if b == 0x9 {
			// do nothing
		} else if b == 0xC {
			screen.Clear()
		} else if b == 0x11 {
			exitFunc()
			return
		} else if b != 0 {
			if err = commsClient.Write([]byte{b}); err != nil {
				status.Text = "Offline"
				status.Refresh()
				displayError(err, w)
			}
		}
	})

	w.ShowAndRun()
}

func execAction(cmd string, args ...string) (string, error) {
	var (
		out []byte
		err error
	)

	if out, err = exec.Command(cmd, args...).Output(); err != nil {
		return "", err
	}

	return string(out), nil
}

func getRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max+1-min) + min
}

func displayError(err error, w fyne.Window) {
	ed := dialog.NewError(err, w)
	ed.Show()
}

func displayMessage(title string, msg string, w fyne.Window) {
	ed := dialog.NewInformation(title, msg, w)
	ed.Show()
}

func saveConnectionFilesToConfig() error {

	var (
		err error
	)

	if configDir, err = CreateConfigDirectory(); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(configDir, "telstar.yml"), []byte(telstarConnection), 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(configDir, "nxtel.yml"), []byte(nxtelConnection), 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(configDir, "ccl4.yml"), []byte(ccl4Connection), 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(configDir, "teefax.yml"), []byte(teefaxConnection), 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(configDir, "eol-bbs.yml"), []byte(eolBbsConnection), 0644); err != nil {
		return err
	}

	return nil
}

func SaveEndpoint(filename string, endpoint comms.Endpoint) error {

	var (
		data []byte
		err  error
	)

	if data, err = yaml.Marshal(endpoint); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil

}

// CreateConfigDirectory creates a config directory called telstar in the platform specific place.
func CreateConfigDirectory() (string, error) {

	var (
		dirName       string
		configDirName string
		err           error
	)

	if dirName, err = os.UserConfigDir(); err != nil {
		return "", err
	}
	configDirName = path.Join(dirName, "telstar/telstar-client")

	if err = os.MkdirAll(configDirName, fs.ModeDir+0777); err != nil {
		return "", err
	}

	return configDirName, nil
}

func setFullScreen(w fyne.Window, waitSeconds int) {

	time.Sleep(time.Second * time.Duration(waitSeconds))
	w.SetFullScreen(true)

}
