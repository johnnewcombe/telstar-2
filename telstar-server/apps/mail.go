package apps

import (
	"github.com/johnnewcombe/telstar/config"
)

func MailSend(sessionId string, settings config.Config, args []string) (bool, error) {
	// maybe internal/external/telegram etc.
	return true, nil
}

func MailList(sessionId string, settings config.Config, args []string) {

}

func MailGet(sessionId string, settings config.Config, args []string) {

}
