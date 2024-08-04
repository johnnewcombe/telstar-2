package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/utils"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/*
All functions return error where more than one thing can go wrong, where only one thing can go wrong,
ok (bool) is returned.

All functions return the natural errors to main() annotated only if appropriate. Helpful
messages are returned to the user.

Where errors are to be created, use errors.New() or fmt.Errorf() as appropriate.

*/

const (
	VERSION     = "2.0"
	REGEXPAGEID = "^[0-9]+[a-z]$"
	REGEXUSERID = "^[0-9]+$"
)

// minimal frame type to use in frame validation
type Frame struct {
	PID       Pid    `json:"pid" bson:"pid"`
	FrameType string `json:"frame-type" bson:"frame-type"`
}

func (f *Frame) IsValid() bool {
	return len(f.PID.FrameId) == 1
}

type AddFramesResponse struct {
	RecordsAdded   int `json:"records added"`
	LastFrameAdded struct {
		PageNo  int    `json:"page-no"`
		FrameId string `json:"frame-id"`
	} `json:"last-frame-added"`
}

type DeleteFrameResponse struct {
	Result string `json:"result"`
}

type Pid struct {
	PageNumber int    `json:"page-no" bson:"page-no"`
	FrameId    string `json:"frame-id" bson:"frame-id"`
}
type EditTfFrame struct {
	PID          Pid     `json:"pid" bson:"pid"`
	Visible      bool    `json:"visible" bson:"visible"`
	Content      Content `json:"content" bson:"content"`
	RoutingTable []int   `json:"routing-table" bson:"routing-table"`
	AuthorId     string  `json:"author-id" bson:"author-id"`
	StaticPage   bool    `json:"static-page" bson:"static-page"`
	FrameType    string  `json:"frame-type" bson:"frame-type"`
}
type Content struct {
	Data string `json:"data" bson:"data"`
	Type string `json:"type" bson:"type"`
}

func main() {

	logger.LogInfo.Printf("Telstar Viewdata Server v%s Bulk Upload Utility. Copyright John Newcombe (2021)\n", VERSION)

	var (
		flags           *flag.FlagSet
		user            string
		password        string
		apiUrl          string
		sourceDirectory string
		primary         bool
		purge           bool
		basePage        int
		includeUnsafe   bool
		err             error
		respData        ResponseData
		token           string
		//pageId            string
		//pageNo            int
		//frameId           rune
		addFramesResponse AddFramesResponse
		//deleteFrameResponse DeleteFrameResponse
	)

	// FIXME add support for zxnet files stored at a url i.e. make source directory support urls
	flags = flag.NewFlagSet("telstar-upload", flag.ExitOnError)
	flags.StringVar(&user, "user", "2222222222", "Username used to access the Telstar API.")
	flags.StringVar(&password, "password", "1234", "Password used to access the Telstar API..")
	flags.StringVar(&apiUrl, "api-url", "localhost:25233", "Url of the Telstar API..")
	flags.StringVar(&sourceDirectory, "source-directory", "./data/frames", "Path to the directory containing the frames (.json) to be uploaded..")
	flags.BoolVar(&primary, "primary", false, "Indicates that the Primary database should be used.")
	flags.IntVar(&basePage, "base-page", 0, "Limits updates to this base page and above.")
	flags.BoolVar(&purge, "purge", true, "Indicates that all follow on frames should be deleted.")
	flags.BoolVar(&includeUnsafe, "include-unsafe", false, "Ensures that all files are imported including unsafe frame types e.g. 'response' and 'gateway'.")

	if err = flags.Parse(os.Args[1:]); err != nil {
		logger.LogError.Printf("%v%", err)
		os.Exit(1)
	}

	// remove data and time and set a prefix
	//log.SetFlags(0)
	//log.SetPrefix("")
	//log.SetOutput(os.Stdout)

	if !strings.HasPrefix(apiUrl, "http") {
		apiUrl = "http://" + apiUrl
	}

	// Login
	if respData, err = cmdLogin(apiUrl, user, password); err != nil {
		log.Fatal(err)
	}
	token = respData.Token

	logger.LogInfo.Println("Logged in successfully.")

	// add frames
	if respData, err = cmdAddFrames(apiUrl, sourceDirectory, primary, token, basePage, includeUnsafe); err != nil {
		log.Fatal(err)
	}

	// sort out frameid to purge from
	if err = json.Unmarshal([]byte(respData.Body), &addFramesResponse); err != nil {
		log.Fatal(err)
	}

	logger.LogInfo.Printf("Frames added: %d.", addFramesResponse.RecordsAdded)

	// Purging:
	// We can keep track of uploaded frames PageNo >= 100 in the 'framesAdded' map
	// and determine the highest frame for each page and delete with 'purge' from
	// there. Purge can only sensibly be used when uploading rss frames where we can
	// guarantee that higher numbered frames that exist in the data base are no longer
	// needed.

	if purge {

		// the framesAdded map has all of the frames >= 100 that were added
		for key := range framesAdded {

			var pageId string

			//get each pageId from the map and see if the follow on frame exists. If not
			// then it wasn't added, therefore we can delete (purge) from this point onwards.
			pageId, err = utils.GetFollowOnPageId(key)

			if !framesAdded[pageId] {
				// delete frame
				logger.LogInfo.Printf("Purging frames from: %s.\r\n", pageId)

				if respData, err = cmdDeleteFrame(apiUrl, pageId, primary, token); err != nil {
					logger.LogInfo.Printf("Unable to purge frames from %s, frames may not exist.\r\n", pageId)
				} else {
					logger.LogInfo.Printf("Frames purged from: %s.\r\n", pageId)
				}
			}
		}
	}
}

func isValidPageId(pageId string) bool {

	// TODO should this be pre-compiled rather than each time the method is called?
	regExFrame, err := regexp.Compile(REGEXPAGEID)
	if err != nil {
		return false
	}
	return regExFrame.MatchString(pageId)
}

func isValidUserId(userId string) bool {

	// TODO should this be pre-compiled rather than each time the method is called?
	regExFrame, err := regexp.Compile(REGEXUSERID)
	if err != nil {
		return false
	}
	return regExFrame.MatchString(userId)
}

func parseFrame(jsonFrame string) (Frame, error) {

	var frame Frame
	frameBytes := []byte(jsonFrame)

	if !json.Valid(frameBytes) {
		return frame, fmt.Errorf("validating frame: invalid json")
	}

	if err := json.Unmarshal(frameBytes, &frame); err != nil {
		return frame, fmt.Errorf("parsing json: invalid")
	}

	return frame, nil
}

func createDefaultRoutingTable(pageNo int) []int {

	var (
		pageNumber float64
	)
	routingTable := make([]int, 11)

	// sort out the entries 0-9 (i.e.keys presses 0-9)
	for n := 0; n < 10; n++ {
		routingTable[n] = n + (pageNo * 10)
	}

	//sort out the hash route
	pageNumber = float64(pageNo)

	for pageNumber > 999 {
		pageNumber = math.Floor(pageNumber / 10)
	}
	routingTable[10] = int(pageNumber)

	return routingTable
}

func checkArgCount(args []string, argCount int) (ok bool) {

	args = deleteEmpty(args)
	lastarg := strings.ToLower(args[len(args)-1])

	//TODO check for primary or secondary and set count accordingly
	if lastarg == globals.DBPRIMARY || lastarg == globals.DBSECONDARY {
		argCount += 1
	}

	return len(args) == argCount
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
func getPidFromFileName(filename string) (pid Pid, ok bool) {

	var (
		pgNo   int
		err    error
		extLen int
	)

	_, filename = filepath.Split(filename)

	//see https://regex101.com/
	var regexEditTf = regexp.MustCompile("^[0-9]{1,10}[a-z].edit.tf$")
	var regexZxNet = regexp.MustCompile("^[0-9]{1,10}[a-z].zxnet$")

	if regexEditTf.MatchString(filename) || regexZxNet.MatchString(filename) {

		extLen = len(path.Ext(filename))

		// get the len of the strings to remove to get the
		if strings.ToLower(path.Ext(filename)) == ".tf" {
			extLen = 8 // special case due to two dots in filename
		} else {
			extLen = len(path.Ext(filename))
		}

		pns := filename[0 : len(filename)-(extLen+1)]
		pnsl := len(pns)
		if pgNo, err = strconv.Atoi(pns); err != nil {
			return pid, false
		}
		pid.PageNumber = pgNo
		pid.FrameId = filename[pnsl : pnsl+1]

		return pid, true

	}

	return pid, false // why is this needed?? after a log.Fatal
}

func exitWithHelp() {

	displayCopyright()

	fmt.Println("  Usage: telstar-util <command> <url> [args...]")
	fmt.Println("")
	fmt.Println("  e.g.   telstar-util login        <url> <user id> <password>")
	fmt.Println("         telstar-util getframe     <url> <page id> [primary|secondary]")
	fmt.Println("         telstar-util getframes    <url> <directory name> [primary|secondary]")
	fmt.Println("         telstar-util addframe     <url> <file name> [primary|secondary]")
	fmt.Println("         telstar-util addframes    <url> <directory name> [primary|secondary]")
	fmt.Println("         telstar-util deleteframe  <url> <page id> [primary|secondary]")
	fmt.Println("         telstar-util purgeframes  <url> <page id> [primary|secondary]")
	fmt.Println("         telstar-util adduser      <url> <user id> <password>")
	fmt.Println("         telstar-util deleteuser   <url> <user id>")
	fmt.Println("         telstar-util getstatus    <url>")
	fmt.Println("")
	fmt.Println("All commands except 'getstatus' require the login command to be issued") // 70
	fmt.Println("first. If the login command is successful a token will be stored in")
	fmt.Println("the following file:")
	fmt.Println("")
	fmt.Println("         telstar-util.tok")
	fmt.Println("")
	fmt.Println("This file will be used in subsequent requests to the API. For this to work,")
	fmt.Println("'telstar-util' must have write permissions to the working directory.")
	fmt.Println("")
	//	fmt.Println("The special command 'edittf' will convert a file containing an edit.tf url")
	//	fmt.Println("into a json frame and save to '<filename>.json'.")
	//	fmt.Println("")

	os.Exit(0)
}

func displayCopyright() {
	fmt.Println("")
	fmt.Println("Telstar Utility v2.0 (c) 2021 John Newcombe")
	fmt.Println("")
}
