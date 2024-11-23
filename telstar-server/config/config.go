package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// The defaults here should line up with the default config. Both are needed in case a the user removes
	// a line from the config etc.
	Api struct {
		Host string `yaml:"host" env:"TELSTAR_API_HOST" env-default:"0.0.0.0"`
	} `yaml:"api"`
	Pad struct {
		Host string `yaml:"host" env:"TELSTAR_PAD_HOST" env-default:"0.0.0.0"`
		DLE  byte   `yaml:"dle" env:"TELSTAR_PAD_DLE" env-default:"0x10"`
	} `yaml:"pad"`

	Server struct {
		Host                    string `yaml:"host" env:"TELSTAR_SERVER_HOST" env-default:"0.0.0.0"`
		DisplayName             string `yaml:"display_name" env:"TELSTAR_SERVER_DISPLAY_NAME" env-default:"DUKE"`
		HidePageId              bool   `yaml:"hide_page_id" env:"TELSTAR_HIDE_PAGE_ID" env-default:"FALSE"`
		HideCost                bool   `yaml:"hide_cost" env:"TELSTAR_HIDE_COST" env-default:"FALSE"`
		DisableVerticalRollOver bool   `yaml:"disable_vertical_rollover" env:"TELSTAR_DISABLE_VERTICAL_ROLLOVER" env-default:"TRUE"`
		EditTfTitleRows         int    `yaml:"edittf_title_rows" env:"TELSTAR_EDITTF_TITLE_ROWS"  env-default:"4"`
		CarouselDelay           int    `yaml:"carousel_delay" env:"TELSTAR_CAROUSEL_DELAY"  env-default:"16"`
		Antiope                 bool   `yaml:"antiope" env:"TELSTAR_ANTIOPE" env-default:"FALSE"`
		DLE                     byte   `yaml:"dle" env:"TELSTAR_SERVER_DLE" env-default:"0x10"`
		Pages                   struct {
			StartPage         int `yaml:"start_page" env:"TELSTAR_START_PAGE"  env-default:"99"`
			LoginPage         int `yaml:"login_page" env:"TELSTAR_LOGIN_PAGE"  env-default:"990"`
			MainIndexPage     int `yaml:"main_index_page" env:"TELSTAR_MAIN_INDEX_PAGE"  env-default:"0"`
			ResponseErrorPage int `yaml:"response_error_page" env:"TELSTAR_RESPONSE_ERROR_PAGE" env-default:"9902"`
			GatewayErrorPage  int `yaml:"gateway_error_page" env:"TELSTAR_GATEWAY_ERROR_PAGE" env-default:"9903"`
		} `yaml:"pages"`

		Authentication struct {
			Required bool `yaml:"required" env:"TELSTAR_REQUIRES_AUTHENTICATION"  env-default:"FALSE"`
			//RootUserId int  `yaml:"root_user_id" env:"TELSTAR_ROOT_USER_ID"  env-default:"0"`
		} `yaml:"authentication"`

		Strings struct {
			DefaultNavMessage          string `yaml:"default_nav_message" env:"TELSTAR_DEFAULT_NAV_MESSAGE"  env-default:"[B][n][Y]Select item or[W]*page# : [_+]"`
			DefaultPageNotFoundMessage string `yaml:"default_page_not_found_message" env:"TELSTAR_DEFAULT_PAGE_NOT_FOUND_MESSAGE"  env-default:"[B][n][Y]Page not Found :[W]"`
			DefaultHeaderText          string `yaml:"default_header_text" env:"TELSTAR_DEFAULT_HEADER_TEXT" env-default:"[G]T[R]E[C]L[B]S[W]T[M]A[Y]R"`
		} `yaml:"strings"`
	} `yaml:"server"`

	Database struct {
		Connection string `yaml:"connection" env:"TELSTAR_DBCON"  env-default:"mongodb://mongoadmin:secret@telstar-mongo:27017"`
		Collection string `yaml:"collection" env:"TELSTAR_DBCOLLECTION"  env-default:"PRIMARY"`
	} `yaml:"database"`

	General struct {
		Parity          bool   `yaml:"parity" env:"TELSTAR_PARITY"  env-default:"FALSE"`
		VolumeDirectory string `yaml:"volume_directory" env:"TELSTAR_VOLUME_DIRECTORY"  env-default:"/opt/telstar/volume/"`
	} `yaml:"general"`
}

func GetConfig(filename string) (Config, error) {

	var (
		cfg Config
		err error
	)

	// get config from file if the filename is specified
	// otherwise uses env vars
	if len(filename) > 0 {
		if cfg, err = readConfig(filename); err != nil {
			if cfg, err = readEnv(); err != nil {
				return cfg, err
			}
		}

	} else {
		if cfg, err = readEnv(); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

// readEnv just uses env for settings, no configuration file
func readEnv() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// readConfig reads the config and overwrites these settings if the env var is set.
// If no setting in config then the env var default is used.
func readConfig(filename string) (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
