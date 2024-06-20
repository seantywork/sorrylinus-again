package controller

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SOLIAGAIN_CONFIG struct {
	Debug       bool   `yaml:"debug"`
	ExternalUrl string `yaml:"externalUrl"`
	ServeAddr   string `yaml:"serveAddr"`
	ServePort   int    `yaml:"servePort"`
	MaxFileSize int64  `yaml:"maxFileSize"`
	Stream      struct {
		TurnServerAddr []struct {
			Addr string `yaml:"addr"`
			Id   string `yaml:"id"`
			Pw   string `yaml:"pw"`
		} `yaml:"turnServerAddr"`
		RtcpPLIInterval        int      `yaml:"rtcpPLIInterval"`
		UploadDest             string   `yaml:"uploadDest"`
		ExtAllowList           []string `yaml:"extAllowList"`
		UdpBufferByteSize      int      `yaml:"udpBufferByteSize"`
		SignalPort             int      `yaml:"signalPort"`
		SignalPortExternal     int      `yaml:"signalPortExternal"`
		RtpReceivePort         int      `yaml:"rtpReceivePort"`
		RtpReceivePortExternal int      `yaml:"rtpReceivePortExternal"`
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
