package controller

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SOLIAGAIN_CONFIG struct {
	ExternalUrl string `yaml:"externalUrl"`
	ServeAddr   string `yaml:"serveAddr"`
	ServePort   string `yaml:"servePort"`
	MaxFileSize int64  `yaml:"maxFileSize"`
	Stream      struct {
		TurnServerAddr    string   `yaml:"turnServerAddr"`
		RtcpPLIInterval   int      `yaml:"rtcpPLIInterval"`
		UploadDest        string   `yaml:"uploadDest"`
		ExtAllowList      []string `yaml:"extAllowList"`
		UdpBufferByteSize int      `yaml:"udpBufferByteSize"`
		SignalPort        string   `yaml:"signalPort"`
		RtpReceivePort    string   `yaml:"rtpReceivePort"`
	} `yaml:"stream"`
	Utils struct {
		UseCompress bool `yaml:"useCompress"`
	} `yaml:"utils"`
}

var CONF *SOLIAGAIN_CONFIG

func LoadConfig() error {

	soliagain_conf := SOLIAGAIN_CONFIG{}

	file_b, err := os.ReadFile("./config.yaml")

	if err != nil {

		return err
	}

	err = yaml.Unmarshal(file_b, &soliagain_conf)

	if err != nil {

		return err
	}

	CONF = &soliagain_conf

	return nil

}
