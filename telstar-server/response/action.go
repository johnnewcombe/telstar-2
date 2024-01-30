package response

import (
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar/apps"
	"bitbucket.org/johnnewcombe/telstar/config"
	"context"
	"os/exec"
	"time"
)

func action(sessionId string, frame *types.Frame, args []string, settings config.Config) (string, error) {

	var (
		authenticated bool
		ok            bool
		cmd           string
		cmdOut        string
		err           error
	)
	cmd = frame.ResponseData.Action.Exec
	if len(cmd) >0 {
		switch cmd {
		case "telstar.login":

			if authenticated, err = apps.Login(sessionId, settings, args); err != nil {
				logger.LogError.Println(err)
			}
			// we update the post action frame depending upon the success of login
			if authenticated {
				logger.LogInfo.Printf("User %s successfully logged in.\r\n", args[0])
			} else {
				logger.LogWarn.Printf("Unsuccessfully login attempt for user %s.\r\n", args[0])
				// TODO set the login failed error page to be next shown
				// we update the post action frame depending upon the success of login
				//frame.ResponseData.Action.PostActionFrame.PageNumber = settings.Pages.LoginPage
			}

		case "telstar.mail.send":
			if ok, err = apps.MailSend(sessionId, settings, args); err != nil {
				logger.LogError.Println(err)
			}
			if ok {
				logger.LogInfo.Printf("User %s sent mail message.\r\n", args[0])
			} else {
				// TODO set the login failed error page to be next shown
			}

		case "telstar.chat.send":
			if ok, err = apps.ChatSend(sessionId, settings, args); err != nil {
				logger.LogError.Println(err)
			}
			if ok {
				logger.LogInfo.Printf("User %s sent chat message.\r\n", args[0])
			} else {
				// TODO set the login failed error page to be next shown
			}

		default:
			logger.LogInfo.Printf("Executing %s %s.\r\n", cmd, args)
			if cmdOut, err = execActionWithCtx(cmd, args...); err != nil {
				return cmdOut, err
			}
		}
		//else {
		//	// we will get here if the exec parameter was empty i.e. don't execute anything
		//	// in which case, instead of raising an error we will simply return a redirect to the
		//}
	}
	return cmdOut, nil
}

func execActionWithCtx(cmd string, args ...string) (string, error) {
	var (
		out []byte
		err error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if out, err = exec.CommandContext(ctx, cmd, args...).CombinedOutput(); err != nil {
		return string(out), err
	}
	return string(out), nil
}
