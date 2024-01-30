package pad

import (
	"fmt"
	"io/ioutil"
	"strings"

	"bitbucket.org/johnnewcombe/telstar-library/logger"
)

type Host struct {
	name     string
	endpoint string
}

func GetHosts(filename string) map[string]string {

	txt, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.LogError.Print(err)
	}
	return parseHosts(string(txt))
}

func parseHosts(text string) map[string]string {

	//var hosts []Port
	var hosts = make(map[string]string)
	//var host Port

	//replace cr fo nl irrespective of platform
	text = strings.Replace(text, "\r", "\n", -1)
	lines := strings.Split(text, "\n")
	lines = cleanStringSlice(lines)

	for _, line := range lines {

		// replace tab
		line = strings.Replace(line, "\t", " ", -1)
		elements := strings.Split(line, " ")
		elements = cleanStringSlice(elements)

		//add a host type to a list based on the following
		if len(elements) == 2 {
			hosts[elements[0]] = elements[1]
		}
	}

	return hosts

}

func cleanStringSlice(slice []string) []string {

	var result []string
	for _, str := range slice {
		if len(str) > 0 {
			result = append(result, str)
		}
	}
	return result
}

func formatHosts(hosts map[string]string) string {

	var msg = "\x0cHOSTS:\r\n\r\n"
	for k, _ := range hosts {
		msg += fmt.Sprintf(" CALL %s\r\n", k)
	}
	msg += "\r\nSee HELP CALL for details.\r\n\r\n"

	return msg
}
