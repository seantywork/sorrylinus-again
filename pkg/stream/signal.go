package stream

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	ch "github.com/seantywork/sorrylinus-again/pkg/com"
)

type SIGNAL_INFO struct {
	Command string `json:"command"`
	Status  string `json:"status"`
	Data    string `json:"data"`
}

var UPGRADER = websocket.Upgrader{}

var roomPeerConnections = make(map[string][]peerConnectionState)

var roomTrackLocals = make(map[string]map[string]*webrtc.TrackLocalStaticRTP)

type peerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	websocket      *ch.ThreadSafeWriter
}

func SignalDispatcher() {

	for range time.NewTicker(time.Millisecond * 10).C {

		for k, _ := range roomPeerConnections {

			dispatchKeyFrame(k)
		}

	}
}
