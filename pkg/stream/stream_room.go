package stream

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
)

var addr string = ":8080"

// lock for peerConnections and trackLocals
var listLock sync.RWMutex
var peerConnections []peerConnectionState
var trackLocals map[string]*webrtc.TrackLocalStaticRTP

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

func StartStreamServerForRoom() {

	go createRoomSignalHandlerForWS()

	router := CreateGenericServer()

	trackLocals = map[string]*webrtc.TrackLocalStaticRTP{}

	router.GET("/", func(c *gin.Context) {

		c.HTML(200, "room.html", gin.H{
			"title": "Room",
		})

	})

	go func() {
		for range time.NewTicker(time.Second * 3).C {
			dispatchKeyFrame()
		}
	}()

	router.Run(addr)

}

func CreateStreamServerForRoom() (*gin.Engine, error) {

	go createRoomSignalHandlerForWS()

	router := CreateGenericServer()

	trackLocals = map[string]*webrtc.TrackLocalStaticRTP{}

	router.GET("/", func(c *gin.Context) {

		c.HTML(200, "room.html", gin.H{
			"title": "Room",
		})

	})

	go func() {
		for range time.NewTicker(time.Second * 3).C {
			dispatchKeyFrame()
		}
	}()

	return router, nil
}

// dispatchKeyFrame sends a keyframe to all PeerConnections, used everytime a new user joins the call
func dispatchKeyFrame() {
	listLock.Lock()
	defer listLock.Unlock()

	for i := range peerConnections {
		for _, receiver := range peerConnections[i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = peerConnections[i].peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
