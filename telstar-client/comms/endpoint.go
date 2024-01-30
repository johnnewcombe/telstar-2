package comms

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

const DEFAULT_CONNECTION = ".recent.yml"

//TODO move this to config package??

type Endpoint struct {
	Name    string `yaml:"name"`
	Address struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"address"`

	Serial struct {
		Port      string `yaml:"port"`
		Baud      int    `yaml:"baud"`
		Parity    bool   `yaml:"parity"`
		ModemInit string `yaml:"modeminit"`
	} `yaml:"serial"`

	Init struct {
		Telnet    bool   `yaml:"telnet"`
		InitChars []byte `yaml:"initchars"`
	} `yaml:"init"`
}

func (e *Endpoint) IsSerial() bool {

	if len(e.Serial.Port) > 0 && e.Serial.Baud > 0 {
		return true
	}
	return false
}

func GetDefaultEndPoint() Endpoint {

	ep := Endpoint{Name: "Telstar"}
	ep.Address.Host = "glasstty.com"
	ep.Address.Port = 6502
	ep.Init.Telnet = false
	ep.Init.InitChars = []byte{}

	return ep
}

func GetEndpoint(filename string) (Endpoint, error) {

	var ep Endpoint
	err := cleanenv.ReadConfig(filename, &ep)
	if err != nil {
		return ep, err
	}
	if len(ep.Name) == 0 {
		return ep, fmt.Errorf("%s is invalid", filename)
	}

	return ep, nil
}




