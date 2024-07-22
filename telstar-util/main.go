package main

import (
	"bitbucket.org/johnnewcombe/telstar-util/cmd"
	"os"
)

//see https://regex101.com/

func main() {
	cmd.Execute()
	os.Exit(0)
}

/*
func main_old() {

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
*/
