package main

/* NOTE NOTE
These tests are dependent on a vanila system existing with API listening on port 25233
These tests test the compiled version of telstar-util so this must be compiled first
Some tests will only work if a valid login has occured
*/

// TODO check dest folder is empty before the getframes test and check it is not empty at the end of the test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

const (
	goodurl  = "localhost:25233"
	badurl   = "localhost:99"
	userId   = "0"
	password = "telstarapisecret"
)

// this is run whenever a test is executed to provide setup
// and teardown options
func TestMain(m *testing.M) {
	var (
		sout string
		err  error
	)

	// call flag.Parse() here if TestMain uses flags
	log.Println("building telstar-util")
	if sout, err = preTestBuild(); err != nil {
		fmt.Println(err)
	}
	fmt.Print(sout)

	// only needed in order to call addframe
	if sout, err = preTestLogin(goodurl, userId, password); err != nil {
		fmt.Println(err)
	}
	fmt.Print(sout)

	if sout, err = preTestAddFrame(goodurl, "./testframes/0a.json"); err != nil {
		fmt.Println(err)
	}
	fmt.Print(sout)

	os.Exit(m.Run())
}

func Test_deleteUser(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)
	_, _ = preTestAddUser(goodurl, "100", "password")

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputUserId  string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"deleteuser, userid 100", goodurl, "100", "telstar-util: user deleted\n", false},
		{"deleteuser, userid not exist", goodurl, "100", "telstar-util: http error: bad request\n", true},
		{"deleteuser, bad userid", goodurl, "a", "telstar-util: http error: bad request\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"deleteuser", test.inputUrl, test.inputUserId}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_deleteUser_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputUserId  string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"deleteuser, not logged in", goodurl, "0", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"deleteuser", test.inputUrl, test.inputUserId}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_addUser(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	//url frameid
	type Test struct {
		description   string
		inputUrl      string
		inputUserId   string
		inputPassword string
		wantResponse  string
		wantErr       bool
	}

	var tests = []Test{
		{"adduser, userid 0", goodurl, "0", "telstarapisecret", "telstar-util: user added\n", false},
		{"adduser, bad userid", goodurl, "a", "telstarapisecret", "telstar-util: http error: bad request\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"adduser", test.inputUrl, test.inputUserId, test.inputPassword}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_addUser_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description   string
		inputUrl      string
		inputUserId   string
		inputPassword string
		wantResponse  string
		wantErr       bool
	}

	var tests = []Test{
		{"adduser, not logged in", goodurl, "0", "telstarapisecret", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"adduser", test.inputUrl, test.inputUserId, test.inputPassword}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_publishFrame(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputFrameId string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"publishframe, frame exists", goodurl, "0a", "telstar-util: frame published\n", false},
		{"publishframe, frame not exist", goodurl, "10000a", "telstar-util: http error: bad request\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"publishframe", test.inputUrl, test.inputFrameId}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_publishFrame_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputFrameId string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"publishframe, not logged in", goodurl, "0a", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"publishframe", test.inputUrl, test.inputFrameId}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_deleteFrame(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputFrameId string
		inputPrimary string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"deleteframe, frame exists", goodurl, "0a", "", "telstar-util: frame deleted\n", false},
		{"deleteframe, frame exists, primary", goodurl, "0a", "primary", "telstar-util: frame deleted\n", false},
		{"deleteframe, frame not exist", goodurl, "0a", "", "telstar-util: http error: bad request\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"deleteframe", test.inputUrl, test.inputFrameId, test.inputPrimary}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_deleteFrame_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputFrameId string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"deleteframe, not loggedin", goodurl, "0a", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"deleteframe", test.inputUrl, test.inputFrameId}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_addFrames(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	type Test struct {
		description    string
		inputUrl       string
		inputDirectory string
		inputPrimary   string
		wantResponse   string
		wantErr        bool
	}
	var tests = []Test{
		{"addframes, directory does not exist", goodurl, "testframes/baddirectory", "", "telstar-util: open testframes/baddirectory: no such file or directory\n", true},
		{"addframes", goodurl, "testframes/upload", "", "telstar-util: 22 frames updated\n", false},
		{"addframes", goodurl, "testframes/upload", "primary", "telstar-util: 22 frames updated\n", false},
	}
	// run tests
	for _, test := range tests {

		var args = []string{"addframes", test.inputUrl, test.inputDirectory, test.inputPrimary}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_addFrames_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description    string
		inputUrl       string
		inputDirectory string
		wantResponse   string
		wantErr        bool
	}

	var tests = []Test{
		{"addframes, not loggedin", goodurl, "testframes/download", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"addframes", test.inputUrl, test.inputDirectory}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_addFrame(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	type Test struct {
		description  string
		inputUrl     string
		inputFile    string
		inputPrimary string
		wantResponse string
		wantErr      bool
	}
	var tests = []Test{
		{"addframe, logged in, file not exist", goodurl, "testframes/1000a.json", "", "telstar-util: open testframes/1000a.json: no such file or directory\n", true},
		{"addframe, logged in, bad file type", goodurl, "testframes/0a.jsn", "", "telstar-util: open testframes/0a.jsn: no such file or directory\n", true},
		{"addframe, logged in", goodurl, "testframes/0a.json", "", "telstar-util: frame updated\n", false},
		{"addframe, logged in, primary", goodurl, "testframes/0a.json", "primary", "telstar-util: frame updated\n", false},
	}
	// run tests
	for _, test := range tests {

		var args = []string{"addframe", test.inputUrl, test.inputFile, test.inputPrimary}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_addFrame_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description    string
		inputUrl       string
		inputFrameFile string

		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"addframe, not loggedin", goodurl, "testframes/0a.json", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"addframe", test.inputUrl, test.inputFrameFile}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_getframes(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	//url frameid
	type Test struct {
		description        string
		inputUrl           string
		inputSaveDirectory string
		inputPrimary       string
		wantResponse       string
		wantErr            bool
	}

	var tests = []Test{
		{"getframes, logged in, bad save directory", goodurl, "testframes/badfolder", "", "telstar-util: open testframes/badfolder/3001a.json: no such file or directory\n", true},
		{"getframes, logged in, good frameid", goodurl, "testframes/download", "", "", false},
		{"getframes, logged in, good frameid, primary", goodurl, "testframes/download", "primary", "", false},
		//TODO check dest folder is empty before the run and check it is not empty at the end of a run
	}

	// run tests
	for _, test := range tests {

		var args = []string{"getframes", test.inputUrl, test.inputSaveDirectory}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}
func Test_getFrames_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description        string
		inputUrl           string
		inputSaveDirectory string

		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"getframes, not logged in", goodurl, "testframes/download", "telstar-util: http error: forbidden\n", true},
		//TODO check dest folder is empty before the run and check it is not empty at the end of a run
	}

	// run tests
	for _, test := range tests {

		var args = []string{"getframes", test.inputUrl, test.inputSaveDirectory}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_getframe(t *testing.T) {

	_, _ = preTestLogin(goodurl, userId, password)

	type Test struct {
		description  string
		inputUrl     string
		inputFrameID string
		inputPrimary string
		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		//{"getframe, not loggedin, good frameid, frame exists", goodurl, "0a", "telstar-util: getframe: login required\n", false},
		{"getframe, logged in, bad frameid", goodurl, "0", "", "telstar-util: invalid frame id\n", true},
		{"getframe, logged in, good frameid, frame not exist", goodurl, "1000a", "", "telstar-util: http error: not found\n", true},
		{"getframe, logged in, good frameid, frame exists", goodurl, "0a", "", "{\"pid\": {\"page-no\": 0, \"frame-id\": \"a\"}, \"redirect\": {\"page-no\": 9, \"frame-id\": \"a\"}, \"visible\": true}\n", false},
		{"getframe, logged in, good frameid, frame exists, primary", goodurl, "0a", "primary", "{\"pid\": {\"page-no\": 0, \"frame-id\": \"a\"}, \"redirect\": {\"page-no\": 9, \"frame-id\": \"a\"}, \"visible\": true}\n", false},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"getframe", test.inputUrl, test.inputFrameID, test.inputPrimary}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_getFrame_NoLogin(t *testing.T) {

	_ = preTestLogoff()

	//url frameid
	type Test struct {
		description  string
		inputUrl     string
		inputFrameId string

		wantResponse string
		wantErr      bool
	}

	var tests = []Test{
		{"getframe, not loggedin", goodurl, "0a", "telstar-util: http error: forbidden\n", true},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"getframe", test.inputUrl, test.inputFrameId}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_login(t *testing.T) {

	err := preTestLogoff()

	if err != nil {
		t.Errorf(TEST_ERROR_MESSAGE, "logging off before tests")
	}

	//url userid password
	type Test struct {
		description   string
		inputUrl      string
		inputUserID   string
		inputPassword string
		wantResponse  string
		wantErr       bool
	}

	var tests = []Test{
		{"login good password, bad url", badurl, "0", "telstarapisecret", "telstar-util: Put \"http://localhost:99/login\": dial tcp 127.0.0.1:99: connect: connection refused\n", true},
		{"login bad password", goodurl, "0", "badpassword", "telstar-util: http error: bad request\n", true},
		{"login good password", goodurl, "0", "telstarapisecret", "telstar-util: login successful\n", false},
	}

	// run tests
	for _, test := range tests {

		var args = []string{"login", test.inputUrl, test.inputUserID, test.inputPassword}

		if got, err := runCommand("./telstar-util", args); got != test.wantResponse ||
			hasError(err) != test.wantErr {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func hasError(err error) bool {
	return err != nil
}

func preTestLogin(url string, userId string, password string) (string, error) {

	var (
		err  error
		sout string
	)
	if sout, err = runCommand("./telstar-util", []string{"login", url, userId, password}); err != nil {
		return "", err
	}
	return sout, nil
}

func preTestLogoff() error {

	if err := saveText(TOKENFILE, ""); err != nil {
		return err
	}
	return nil
}

func preTestBuild() (string, error) {
	// ignore errors as there are always warnings with the compile
	sout, err := runCommand("./telstar-util-build", []string{""})

	if err != nil {
		return "", err
	}
	return sout, err
}

func preTestAddUser(url string, userId string, password string) (string, error) {
	var (
		sout string
		err  error
	)
	sout, err = runCommand("./telstar-util", []string{"adduser", url, userId, password})
	if err != nil {
		return "", err
	}
	return sout, nil
}

func preTestAddFrame(url string, filename string) (string, error) {
	var (
		sout string
		err  error
	)
	sout, err = runCommand("./telstar-util", []string{"addframe", url, filename})
	if err != nil {
		return "", err
	}
	sout, err = runCommand("./telstar-util", []string{"addframe", url, filename, "primary"})
	if err != nil {
		return "", err
	}

	return sout, nil
}

func runCommand(name string, args []string) (string, error) {

	//var out bytes.Buffer
	var (
		stdout []uint8
		err    error
	)

	//stdout, err = exec.Command(name, args...).Output() // '...' means 'Unpack Slice'
	cmd := exec.Command(name, args...) // '...' means 'Unpack Slice'
	stdout, err = cmd.Output()

	return string(stdout), err
}
