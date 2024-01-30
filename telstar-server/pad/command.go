package pad

import (
	"regexp"
	"strings"
)

type Command struct {
	name    string
	isValid bool
	arg1    string
	arg2    string
}

var (
	//profileRegEx = regexp.MustCompile("^profile=p[0-9]$")
	profileRegEx  = regexp.MustCompile("(?i)" + "^profile=p[0-9]$")         //(?i) = case insensitive
	paramRegEx    = regexp.MustCompile("(?i)" + "^p[1-9]{1,2}=[0-9]{1,3}$") //(?i) = case insensitive
	cmd           Command
	commandBuffer string
)

func parseCommand(commandLine string) Command {

	var cmd Command
	//var commandList = []string{"help", "profile", "profiles", "exit"}

	commandLine = strings.Trim(strings.ToUpper(commandLine), " ")
	args := strings.Split(commandLine, " ")

	if len(args) > 0 {

		if isParamSetting(args[0]) {
			cmd.name = "P"
		} else if isProfileSetting(args[0]) {
			cmd.name = "PROFILESET"
		} else {
			cmd.name = args[0]
		}

		switch cmd.name {

		case "HOSTS", "PROFILES":
			cmd.isValid = true

		case "HELP":

			if len(args) > 1 {
				//help sub command
				switch args[1] {
				case "PARS", "PROFILE", "CALL":
					cmd.name += args[1]
					cmd.isValid = true
				}
			} else {
				cmd.isValid = true
			}

		case "CALL":

			if len(args) > 1 {
				cmd.isValid = true
				cmd.arg1 = args[1]
			}

		case "PROFILE":
			cmd.isValid = true
			if len(args) > 1 {
				cmd.arg1 = args[1]
			}

		case "PARS":
			cmd.isValid = true

		case "PROFILESET":
			cmd.isValid = true
			cmd.arg1 = parseProfile(args[0])

		case "P":
			cmd.isValid = true
			cmd.arg1, cmd.arg2 = parseParam(args[0])

		}
	}

	return cmd
}

func isStringInSlice(str string, list []string) bool {

	for _, b := range list {
		if b == str {
			return true
		}
	}
	return false
}

func isProfileSetting(commandLine string) bool {

	var result = false
	commandLine = strings.ReplaceAll(commandLine, " ", "")
	//commandLine = strings.Trim(strings.ToLower(commandLine), " ")
	args := strings.Split(commandLine, " ")

	if len(args) > 0 {
		// check it is not a parameter i.e. p2=45
		if profileRegEx.MatchString(args[0]) {
			result = true
		} else {
			result = false
		}
	}
	return result
}
func isParamSetting(commandLine string) bool {

	var result = false
	commandLine = strings.ReplaceAll(commandLine, " ", "")
	args := strings.Split(commandLine, " ")

	if len(args) > 0 {
		// check it is not a parameter i.e. p2=45
		if paramRegEx.MatchString(args[0]) {
			result = true
		} else {
			result = false
		}
	}
	return result
}
func parseParam(paramArg string) (string, string) {

	if isParamSetting(paramArg) {
		paramArg = strings.ReplaceAll(paramArg, " ", "")
		arg := strings.Split(paramArg, "=")
		arg1 := strings.Trim(arg[0], "P")
		arg1 = strings.Trim(arg1, "p")
		arg2 := arg[1]

		return arg1, arg2
	} else {
		return "", ""
	}
}
func parseProfile(profileArg string) string {

	if isProfileSetting(profileArg) {
		profileArg = strings.ReplaceAll(profileArg, " ", "")
		arg := strings.Split(profileArg, "=")
		arg1 := arg[1]
		//arg2 := arg[1]

		return arg1
	} else {
		return ""
	}
}
