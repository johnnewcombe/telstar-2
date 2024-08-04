package globals

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var version string

func GetVersion() (string, error) {

	return strings.Replace(string(version), "\n", "", -1), nil
}
