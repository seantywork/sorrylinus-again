package controller

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SOLIAGAIN_CONFIG struct {
	Debug       bool   `yaml:"debug"`
	ExternalUrl string `yaml:"externalUrl"`
	InternalUrl string `yaml:"internalUrl"`
	ServeAddr   string `yaml:"serveAddr"`
	ServePort   int    `yaml:"servePort"`
	MaxFileSize int64  `yaml:"maxFileSize"`
	TimeoutSec  int    `yaml:"timeoutSec"`
	Com         struct {
		ChannelPort         int `yaml:"channelPort"`
		ChannelPortExternal int `yaml:"channelPortExternal"`
	} `yaml:"com"`
	Sorrylinus struct {
		FrontAddr      string `yaml:"frontAddr"`
		SoliSignalAddr string `yaml:"soliSignalAddr"`
	} `yaml:"sorrylinus"`
	Edition struct {
		ExtAllowList []string `yaml:"extAllowList"`
	} `yaml:"edition"`
	Stream struct {
		TurnServerAddr []struct {
			Addr string `yaml:"addr"`
			Id   string `yaml:"id"`
			Pw   string `yaml:"pw"`
		} `yaml:"turnServerAddr"`
		PeerSignalAttemptSync  int    `yaml:"peerSignalAttemptSync"`
		PeerSignalAddr         string `yaml:"peerSignalAddr"`
		SingleRoomMode         bool   `yaml:"singleRoomMode"`
		RtcpPLIInterval        int    `yaml:"rtcpPLIInterval"`
		UdpBufferByteSize      int    `yaml:"udpBufferByteSize"`
		UdpMuxPort             int    `yaml:"udpMuxPort"`
		UdpEphemeralPortMin    int    `yaml:"udpEphemeralPortMin"`
		UdpEphemeralPortMax    int    `yaml:"udpEphemeralPortMax"`
		RtpReceivePort         int    `yaml:"rtpReceivePort"`
		RtpReceivePortExternal int    `yaml:"rtpReceivePortExternal"`
	} `yaml:"stream"`
	Log struct {
		FlushIntervalSec int `yaml:"flushIntervalSec"`
	} `yaml:"log"`
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
