package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
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
	TOKENFILE   = "telstar-util.tok"
	REGEXPAGEID = "^[0-9]+[a-z]$"
	REGEXUSERID = "^[0-9]+$"
)

//see https://regex101.com/
var regex = regexp.MustCompile("^[0-9]{1,10}[a-z].edit.tf$")
/*
// minimal frame type to use in frame validation
type Frame struct {
	PID Pid `json:"pid" bson:"pid"`
}
func (f *Frame) IsValid() bool {
	return len(f.PID.FrameId) == 1
}
type Pid struct {
	PageNumber int    `json:"page-no" bson:"page-no"`
	FrameId    string `json:"frame-id" bson:"frame-id"`
}
*/


type EditTfFrame struct {
	PID          types.Pid     `json:"pid" bson:"pid"`
	Visible      bool    `json:"visible" bson:"visible"`
	Content      types.Content `json:"content" bson:"content"`
	RoutingTable []int   `json:"routing-table" bson:"routing-table"`
	AuthorId     string  `json:"author-id" bson:"author-id"`
	StaticPage   bool    `json:"static-page" bson:"static-page"`
}




func main() {

	// remove data and time and set a prefix
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(os.Stdout)

	if len(os.Args) < 3 {
		exitWithHelp()
	}

	var (
		respData ResponseData
		token    string
		err      error
	)

	apiUrl := os.Args[2] // url always the first arg after the sub command
	if !strings.HasPrefix(apiUrl, "http") {
		apiUrl = "http://" + apiUrl
	}

	var primary = os.Args[len(os.Args)-1] == "primary" // database always the last arg

	switch strings.ToLower(os.Args[1]) {

	case "getstatus":

		// check args
		if !checkArgCount(os.Args, 3) {
			exitWithHelp()
		}

		//call method
		if respData, err = cmdGetStatus(apiUrl); err != nil {
			log.Fatal(err)
		}

		//output the result
		log.Println(respData.Body)
	case "login":

		// check args
		if !checkArgCount(os.Args, 5) {
			exitWithHelp()
		}

		//call method
		if respData, err = cmdLogin(apiUrl, os.Args[3], os.Args[4]); err != nil {
			log.Fatal(err)
		}

		if err = saveText(TOKENFILE, respData.Token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "getframe":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}

		//load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		//call method
		if respData, err = cmdGetFrame(apiUrl, os.Args[3], primary, token); err != nil {
			log.Fatal(err)
		}

		//output the result
		log.Println(respData.Body)

	case "getframes":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}

		//load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		//call method (args 3 holds directory
		if respData, err = cmdGetFrames(apiUrl, os.Args[3], primary, token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "addframe":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}
		// load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		// call method
		if respData, err = cmdAddFrame(apiUrl, os.Args[3], primary, token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "addframes":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}

		// get token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		//call method
		if respData, err = cmdAddFrames(apiUrl, os.Args[3], primary, token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "deleteframe":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}

		//load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		pageId := os.Args[3]

		if respData, err = cmdDeleteFrame(apiUrl, pageId, primary, false, token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "purgeframes":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}

		//load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		pageId := os.Args[3]

		if respData, err = cmdDeleteFrame(apiUrl, pageId, primary, true, token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "adduser":

		//check args
		if !checkArgCount(os.Args, 5) {
			exitWithHelp()
		}

		// load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		// call method
		if respData, err = cmdAddUser(apiUrl, os.Args[3], os.Args[4], token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	case "deleteuser":

		// check args
		if !checkArgCount(os.Args, 4) {
			exitWithHelp()
		}

		//load token
		if token, err = loadText(TOKENFILE); err != nil {
			log.Fatal(err)
		}

		// call method
		if respData, err = cmdDeleteUser(apiUrl, os.Args[3], token); err != nil {
			log.Fatal(err)
		}

		log.Println(respData.Body)

	default:
		exitWithHelp()
	}
}

func parseFrame(jsonFrame string) (types.Frame, error) {

	var frame types.Frame
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
func getPidFromFileName(filename string) (pid types.Pid, ok bool) {

	var (
		pgNo int
		err  error
	)
	_, filename = filepath.Split(filename)

	if regex.MatchString(filename) {

		pns := filename[0 : len(filename)-9]
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
