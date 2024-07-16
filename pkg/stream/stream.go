package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pion/ice/v3"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"github.com/seantywork/sorrylinus-again/pkg/com"
)

var DEBUG bool

var EXTERNAL_URL string

var INTERNAL_URL string

var RTCP_PLI_INTERVAL time.Duration

var UDP_MUX_PORT int

var UDP_EPHEMERAL_PORT_MIN int

var UDP_EPHEMERAL_PORT_MAX int

var TIMEOUT_SEC int

var TURN_SERVER_ADDR []struct {
	Addr string `json:"addr"`
	Id   string `json:"id"`
	Pw   string `json:"pw"`
}

var api *webrtc.API

func GetTurnServeAddr(c *gin.Context) {

	data_b, err := json.Marshal(TURN_SERVER_ADDR)

	if err != nil {

		c.JSON(http.StatusOK, com.SERVER_RE{Status: "failed", Reply: "error"})
	}

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: string(data_b)})

}

func InitWebRTCApi() {

	settingEngine := webrtc.SettingEngine{}

	mux, err := ice.NewMultiUDPMuxFromPort(UDP_MUX_PORT)

	settingEngine.SetICEUDPMux(mux)
	if err != nil {
		panic(err)
	}

	settingEngine.SetEphemeralUDPPortRange(uint16(UDP_EPHEMERAL_PORT_MIN), uint16(UDP_EPHEMERAL_PORT_MAX))

	api = webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))

}

func recieveTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]*webrtc.TrackLocalStaticRTP,
	peerID string) {

	if _, ok := peerConnectionMap[peerID]; !ok {

		newLocalTrack := webrtc.TrackLocalStaticRTP{}

		peerConnectionMap[peerID] = &newLocalTrack
	}

	localTrack := peerConnectionMap[peerID]

	peerConnection.AddTrack(localTrack)

	fmt.Printf("connection map: %d\n", len(peerConnectionMap))

}

// user is the caller of the method
// if user connects before peer: since user is first, user will create the channel and track and will pass the track to the channel
// if peer connects before user: since peer came already, he created the channel and is listning and waiting for me to create and pass track
func createTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]*webrtc.TrackLocalStaticRTP,
	currentUserID string) {

	if _, err := peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		log.Fatal(err)
	}

	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it

		go func() {
			ticker := time.NewTicker(RTCP_PLI_INTERVAL)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remoteTrack.RtxSSRC())}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		// main track of the broadcaster

		localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, remoteTrack.ID(), remoteTrack.StreamID())
		if newTrackErr != nil {
			log.Fatal(newTrackErr)
		}

		// the channel that will have the local track that is used by the sender
		// the localTrack needs to be fed to the reciever

		if _, ok := peerConnectionMap[currentUserID]; ok {
			// feed the exsiting track from user with this track
			peerConnectionMap[currentUserID] = localTrack
		} else {
			peerConnectionMap[currentUserID] = localTrack
		}

		rtpBuf := make([]byte, 1400)
		for { // for publisher only
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatal(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err := localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				log.Fatal(err)
			}
		}
	})

}

func addTrack(k string, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	com.ListLock.Lock()
	defer func() {
		com.ListLock.Unlock()
		signalPeerConnections(k)
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	//tl := map[string]*webrtc.TrackLocalStaticRTP{
	//	t.ID(): trackLocal,
	//}

	trackLocals[t.ID()] = trackLocal
	return trackLocal
}

func removeTrack(k string, t *webrtc.TrackLocalStaticRTP) {
	com.ListLock.Lock()
	defer func() {
		com.ListLock.Unlock()
		signalPeerConnections(k)
	}()

	delete(trackLocals, t.ID())
}
