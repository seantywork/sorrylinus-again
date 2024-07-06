package stream

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var SIGNAL_ADDR string

var SIGNAL_PORT string

var SIGNAL_PORT_EXTERNAL string

var USER_SIGNAL = make(map[string]*websocket.Conn)

var UPGRADER = websocket.Upgrader{}

type SIGNAL_INFO struct {
	Command string `json:"command"`
	Status  string `json:"status"`
	Data    string `json:"data"`
}

var listLock sync.RWMutex
var peerConnections = make([]peerConnectionState, 0)
var trackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)

type peerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	websocket      *threadSafeWriter
}

type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

func AddSignalHandler(signalPath string, signalHandler func(w http.ResponseWriter, r *http.Request)) {

	http.HandleFunc(signalPath, signalHandler)

}

func StartSignalHandler() {

	go func() {

		for range time.NewTicker(time.Second * 3).C {
			dispatchKeyFrame()
		}

	}()

	signal_addr := SIGNAL_ADDR + ":" + SIGNAL_PORT

	log.Fatal(http.ListenAndServe(signal_addr, nil))

}
