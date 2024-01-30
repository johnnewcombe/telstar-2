package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/file"
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar/api"
	"bitbucket.org/johnnewcombe/telstar/config"
	"bitbucket.org/johnnewcombe/telstar/dal"
	"bitbucket.org/johnnewcombe/telstar/install"
	"bitbucket.org/johnnewcombe/telstar/pad"
	"bitbucket.org/johnnewcombe/telstar/server"
	"flag"
	"fmt"
	"github.com/go-chi/jwtauth"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	CONFIG_FILE        = "telstar.yml"
	VERSION     string = "2.0"

	EnvApiUserid       = "TELSTAR_API_USERID"
	EnvApiPassword     = "TELSTAR_API_PASSWORD"
	DefaultApiUserid   = "2222222222"
	DefaultApiPassword = "1234"
	EnvApiSecret       = "TELSTAR_API_SECRET"
	// TODO this name is for backward compatibility with early versions of Telstar
	EnvCookieSecret = "TELSTAR_COOKIE_SECRET"
)

func main() {

	logger.LogInfo.Printf("Telstar Videotex Server Version %s. Copyright John Newcombe (2021)\n", VERSION)

	var (
		serverPort      int
		padPort         int
		apiPort         int
		init            bool
		initSystemPages bool
		configFile      string
		err             error
	)

	// NOTE: when parsing each sub-command's flag's port will be successfully parsed in each case
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.IntVar(&serverPort, "port", 25232, "TCP port number that the service will listen on")
	serverCmd.BoolVar(&init, "init", false, "Creates/recreates the default system and redirect pages and adds them to the database.")
	serverCmd.BoolVar(&initSystemPages, "init-system-pages", false, "Creates/recreates the default system pages and adds them to the database.")

	padCmd := flag.NewFlagSet("pad", flag.ExitOnError)
	padCmd.IntVar(&padPort, "port", 25233, "TCP port number that the service will listen on")
	//padCmd.BoolVar(&dev, "dev", false, "Loads the telstar-dev.yml configuration file.")

	apiCmd := flag.NewFlagSet("api", flag.ExitOnError)
	apiCmd.IntVar(&apiPort, "port", 25234, "TCP port number that the service will listen on")
	//apiCmd.BoolVar(&dev, "dev", false, "Loads the telstar-dev.yml configuration file.")
	//apiCmd.BoolVar(&init, "init", false, "Creates/recreates the default api system user to the database.")

	if len(os.Args) < 2 {
		fmt.Println(getCmdHelp())
		os.Exit(1)
	}
	switch strings.ToLower(os.Args[1]) {
	case "server":
		if err = serverCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(getCmdHelp())
			logger.LogError.Fatal(err)
		}
	case "pad":
		if err = padCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(getCmdHelp())
			logger.LogError.Fatal(err)
		}
	case "api":
		if err = apiCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(getCmdHelp())
			logger.LogError.Fatal(err)
		}
	default:
		logger.LogError.Println("expected 'server', 'pad' or 'api' subcommands")
		fmt.Println(getCmdHelp())
		os.Exit(1)
	}
	//FIXME sort out response if params are wrong or missing.
	//if !(os.Args[1] == "server" || os.Args[1] == "pad" || os.Args[1] == "api") {
	//	fmt.Println("expected 'server', 'pad' or 'api' subcommands")
	//	os.Exit(1)
	//}

	// note that the config.Config includes definitions for environment variable.
	// get config here so that it can be passed to any sub commands
	var settings config.Config
	wd, _ := os.Getwd()
	configFile = path.Join(wd, CONFIG_FILE)

	if file.Exists(configFile) {

		logger.LogInfo.Printf("Using config file %s.", configFile)
	} else {
		configFile = ""
		logger.LogInfo.Printf("No config file, using environment variables for settings.")
	}

	if settings, err = config.GetConfig(configFile); err != nil {
		logger.LogError.Fatal(err)
	}

	switch os.Args[1] {

	case "server":

		// create the system pages (see below for redirects)
		if init || initSystemPages {

			// create/recreate the system user
			logger.LogInfo.Print("Creating the Manager Account.")
			if err = createSystemManagerAccount(settings); err != nil {
				logger.LogError.Print(err)
			}
			logger.LogInfo.Print("Creating the System Pages.")
			if err := install.CreateSystemPages(settings); err != nil {
				logger.LogError.Print(err)
			}
			logger.LogInfo.Print("Creating the Sample Gateway Pages.")
			if err := install.CreateGatewayPages(settings); err != nil {
				logger.LogError.Print(err)
			}
			logger.LogInfo.Print("Creating the Sample Error Pages.")
			if err := install.CreateErrorPages(settings); err != nil {
				logger.LogError.Print(err)
			}

			// if init then update any redirect pages
			if init {
				logger.LogInfo.Print("Creating the Redirect Pages.")
				if err := install.CreateSystemRedirectPages(settings); err != nil {
					logger.LogError.Print(err)
				}
			}
		}
		// log the server(s) starting
		logger.LogError.Print(server.Start(serverPort, settings))

	case "pad":

		// get the host list from the file "hosts"
		hosts := pad.GetHosts("hosts")
		logger.LogError.Print(pad.Start(padPort, settings, hosts))

	case "api":

		// this is the env for Telstar v2.0 which uses HS256 jwt tokens
		apiSecret := os.Getenv(EnvApiSecret)

		// this is the env for early versions of Telstar which uses a encrypted abstract token
		cookieSecret := os.Getenv(EnvCookieSecret)

		if len(apiSecret) == 0 {
			// support for env vars used in Telstar 0.x and 1.x
			if len(cookieSecret) == 0 {
				logger.LogError.Print("The Telstar API is unable to start as the TELSTAR_API_SECRET environment variable has not been set. \nIdeally, his should be set using at least a 32 character string. The string should be kept secret!")
				return
			} else {
				//The TELSTAR_API_SECRET has not been set. Using TELSTAR_COOKIE_SECRET instead.
				apiSecret = cookieSecret
			}
		}
		api.TokenAuth = jwtauth.New("HS256", []byte(apiSecret), nil)
		// start the api
		logger.LogError.Print(api.Start(apiPort, settings))
		//api.Start(apiPort, settings)
	}
}

func createSystemManagerAccount(settings config.Config) error {

	var (
		apiUserId   string
		apiPassword string
		err         error
	)
	// create/recreate the system user using the env, unless env isn't specified or it happens
	// to match the guest ID, in which case use the default value
	if apiUserId = os.Getenv(EnvApiUserid); len(apiUserId) == 0 || apiUserId == globals.GUEST_USER {
		apiUserId = DefaultApiUserid
	}

	if apiPassword = os.Getenv(EnvApiPassword); len(apiPassword) == 0 {
		apiPassword = DefaultApiPassword
	}

	user := types.User{UserId: apiUserId, Password: apiPassword, Name: "SYSTEM MANAGER", BasePage: 0, Admin: true, ApiAccess: true, Authenticated: false}
	if err = dal.InsertOrReplaceUser(settings.Database.Connection, user); err != nil {
		return err
	}
	return nil
}

/*
func createGuestAccount(settings config.Config) error {

	var (
		user dal.User
		err  error
	)

	user = dal.User{}

	// add the guest user to the database
	user = dal.User{UserId: globals.GUEST_USER, Password: globals.GUEST_PASSWORD, Name: globals.GUEST_NAME}
	if err = dal.InsertOrReplaceUser(settings.Database.Connection, user); err != nil {
		return err
	}
	return nil
}
*/
/*
func createStandardUser(settings config.Config) error {

	var (
		user dal.User
		err  error
	)

	// add the standard user to the database
	user = dal.User{UserId: globals.STD_USER, Password: globals.STD_PASSWORD, Name: "MR GOLD", BasePage: 200}
	if err = dal.InsertOrReplaceUser(settings.Database.Connection, user); err != nil {
		return err
	}

	return nil

}
*/

func getCmdHelp() string {

	//12345678901234567890123456789012345678901234567890123456789012345678901234567890
	return "\nTelstar is a videotex server complete with restful administration api.\nUsage:\n\tgo <command> [arguments]\n\nThe commands are:\n\tserver\t\tstart a telstar videotex server\n\tapi\t\tstart the telstar api server\n\tpad\t\tstart the telstar packet assembler/disassembler server\n"

}

func getOs() string {

	var result string

	switch osysy := runtime.GOOS; osysy {
	case "darwin":
		result = "MacOS"
	case "linux":
		result = "Linux"
	default:
		// freebsd, openbsd,
		// plan9, windows...
		result = osysy
	}

	return result
}
