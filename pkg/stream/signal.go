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

var peerConnections = make([]peerConnectionState, 0)
var trackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)

type peerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	websocket      *ch.ThreadSafeWriter
}

func SignalDispatcher() {

	for range time.NewTicker(time.Second * 3).C {
		dispatchKeyFrame()
	}
}
