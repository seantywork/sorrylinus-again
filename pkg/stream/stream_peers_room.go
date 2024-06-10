package stream

import (
	"fmt"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pion/webrtc/v4"
)

func CreateStreamServerForPeersRoom() (*gin.Engine, error) {

	go createSignalHandlerForWS()

	router := CreateGenericServer()

	peerConnectionMap := make(map[string]chan *webrtc.TrackLocalStaticRTP)

	m := &webrtc.MediaEngine{}

	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "video/VP8", ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        96,
	}, webrtc.RTPCodecTypeVideo); err != nil {

		return nil, err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	/*
		peerConnectionConfig := webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		}
	*/
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:localhost:3478"},
			},
		},
	}
	router.GET("/", func(c *gin.Context) {

		c.HTML(200, "peers_room.html", gin.H{
			"title": "Peers Room",
		})

	})

	router.POST("/peers/room/sdp/m/:meetingId/c/:userID/s/:isSender", func(c *gin.Context) {

		fmt.Println("webrtc room post access")

		isSender, _ := strconv.ParseBool(c.Param("isSender"))

		if isSender {
			fmt.Println("sender")
		} else {

			fmt.Println("receiver")
		}

		userID := c.Param("userID")

		var session Sdp
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		offer := webrtc.SessionDescription{}
		Decode(session.Sdp, &offer)

		// Create a new RTCPeerConnection
		// this is the gist of webrtc, generates and process SDP
		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {

			fmt.Println(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return

		}
		if !isSender {
			recieveTrack(peerConnection, peerConnectionMap, userID)
		} else {
			createTrack(peerConnection, peerConnectionMap, userID)
		}

		peerConnection.SetRemoteDescription(offer)

		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {

			fmt.Println(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, Sdp{Sdp: Encode(answer)})
	})

	return router, nil

}
