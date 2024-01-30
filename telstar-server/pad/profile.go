package pad

import (
	"fmt"
	"strings"
)

const (
	DEFAULT_PROFILE = "P7"
)

type Profile struct {
	name   string
	values []ProfileValue
}
type ProfileValue struct {
	index       int
	value       int
	description string
}

var (
	currentProfile Profile
)
var profileDescriptions = []string{
	"PAD Recall",
	"Echo",
	"Forward",
	"Idle",
	"A. XON/XOFF",
	"Signals",
	"Break",
	"Discard",
	"CR Padding",
	"Line Folding",
	"Speed",
	"Flow Control",
	"LF/CR",
	"LF Padding",
	"Editing",
	"Char Delete",
	"Line Delete",
	"Line Display",
	"Edit Signals",
	"Echo Mask",
	"Parity",
	"Page Wait",
}
var profileP7 = []ProfileValue{ // max desc len 11 chars
	{1, 1, profileDescriptions[0]},
	{2, 1, profileDescriptions[1]},
	{3, 0, profileDescriptions[2]},
	{4, 0, profileDescriptions[3]},
	{5, 0, profileDescriptions[4]},
	{6, 0, profileDescriptions[5]},
	{7, 0, profileDescriptions[6]},
	{8, 0, profileDescriptions[7]},
	{9, 0, profileDescriptions[8]},
	{10, 0, profileDescriptions[9]},
	{11, 255, profileDescriptions[10]},
	{12, 0, profileDescriptions[11]},
	{13, 0, profileDescriptions[12]},
	{14, 0, profileDescriptions[13]},
	{15, 0, profileDescriptions[14]},
	{16, 127, profileDescriptions[15]},
	{17, 24, profileDescriptions[16]},
	{18, 18, profileDescriptions[17]},
	{19, 2, profileDescriptions[18]},
	{20, 64, profileDescriptions[19]},
	{21, 2, profileDescriptions[20]},
	{22, 0, profileDescriptions[21]},
}
var profileP8 = make([]ProfileValue, len(profileP7))

func init(){
	// profileP8 is very similar to profileP7
	copy(profileP8, profileP7)
	profileP8[20] = ProfileValue{21, 0, profileDescriptions[20]}

}
/*
var profileP8 = []ProfileValue{
	{1, 1, "PAD Recall"},
	{2, 1, "Echo"},
	{3, 0, "Forward"},
	{4, 0, "Idle"},
	{5, 0, "Ancillary XON/XOFF"},
	{6, 0, "Signals"},
	{7, 0, "Break"},
	{8, 0, "Discard"},
	{9, 0, "CR Padding"},
	{10, 0, "Line Folding"},
	{11, 255, "Speed"},
	{12, 0, "Terminal Flow Control"},
	{13, 0, "LF/CR"},
	{14, 0, "LF Padding"},
	{15, 0, "Editing"},
	{16, 127, "Char Delete"},
	{17, 24, "Line Delete"},
	{18, 18, "Line Display"},
	{19, 2, "Editing Signals"},
	{20, 64, "Echo Mask"},

	{21, 0, "Parity"}, // this is the only difference from P7

	{22, 0, "Page Wait"},
}
*/
var profiles = []Profile{
	{"P7", profileP7},
	{"P8", profileP8},
}

func formatProfiles(profiles []Profile) string {

	var msg = "\fAVAILABLE PROFILES:\r\n\r\n"
	for _, profile := range profiles {
		msg += fmt.Sprintf("%11s\r\n", profile.name)
	}
	msg += "\r\nPROFILE <NAME> shows profile details.\r\n" + RETURN

	return msg

}

func formatProfile(profile Profile) string {

	var msg string
	msg += "\x0cCURRENT PROFILE:\r\n\r\n"
	msg += fmt.Sprintf("\r\nProfile Name: %s", profile.name)
	msg += formatPars(profile)
	return msg
}

func formatPars(profile Profile) string {

	var msg string
	msg += "\x0cPARAMETERS:\r\n\r\n"
	for i, v := range profile.values {
		if i%2 == 0 {
			msg += fmt.Sprintf("%02d=%03d %-12s", v.index, v.value, v.description)
		} else {
			msg += fmt.Sprintf(" %02d=%03d %s\r\n", v.index, v.value, v.description)
		}
	}
	msg += RETURN

	return msg
}

func getProfile(profileName string) (profile Profile, ok bool) {

	profileName = strings.Trim(strings.ToLower(profileName), " ")
	for _, p := range profiles {
		if strings.ToLower(p.name) == strings.ToLower(profileName) {
			return p, true
		}
	}
	//return Profile{}
	return Profile{}, false
}

func getProfileValue(profile *Profile, index int) (value int, ok bool) {

	if index > len(profile.values) {
		return 0, false
	}
	return profile.values[index-1].value, true
}

func setProfileValue(profile *Profile, index int, value int) bool {

	if index > len(profile.values) {
		return false
	}
	profile.values[index-1].value = value
	return true
}

func setProfile(profileName string) (ok bool) {

	profile, ok := getProfile(profileName)

	// copy profile so that any changes don't affect the standard ones
	// first make sure we have room in the destination
	currentProfile.values = make([]ProfileValue, len(profile.values))

	if ok {
		currentProfile.name = profile.name
		copy(currentProfile.values, profile.values)

		// check for consistency
		if currentProfile.values == nil || len(currentProfile.values) != len(profile.values) {
			return !ok
		}
	}
	return ok
}
