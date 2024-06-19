package stream

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var PEERS_SIGNAL_PATH string

func GetPeersIndex(c *gin.Context) {

	c.HTML(200, "peers.html", gin.H{
		"title": "Peers",
	})

}

func GetPeersSignalAddress(c *gin.Context) {

	s_addr := EXTERNAL_URL + ":" + SIGNAL_PORT + PEERS_SIGNAL_PATH

	c.JSON(http.StatusOK, SERVER_RE{Status: "success", Reply: s_addr})

}

func InitPeersSignalOn(peerSignalPath string) {

	PEERS_SIGNAL_PATH = peerSignalPath

	runPeersSignalHandlerForWS(peerSignalPath)

}

/*
func CreateStreamServerForPeers() (*gin.Engine, error) {

	go createPeersSignalHandlerForWS()

	router := CreateGenericServer()

	peerConnectionMap := make(map[string]*webrtc.TrackLocalStaticRTP)

	m := &webrtc.MediaEngine{}

	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "video/VP8", ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        96,
	}, webrtc.RTPCodecTypeVideo); err != nil {

		return nil, err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{TurnServerAddr},
			},
		},
	}
	router.GET("/", func(c *gin.Context) {

		c.HTML(200, "peers.html", gin.H{
			"title": "Peers",
		})

	})

	router.GET("/peers/room/turn", func(c *gin.Context) {

		c.JSON(http.StatusOK, SERVER_RE{Status: "success", Reply: TurnServerAddr})

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

		var session CLIENT_REQ
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid request"})
			return
		}

		offer := webrtc.SessionDescription{}
		Decode(session.Data, &offer)

		// Create a new RTCPeerConnection
		// this is the gist of webrtc, generates and process SDP
		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {

			fmt.Println(err.Error())

			c.JSON(http.StatusInternalServerError, SERVER_RE{Status: "error", Reply: "failed to process"})

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

			c.JSON(http.StatusInternalServerError, SERVER_RE{Status: "error", Reply: "failed to process"})

			return
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {

			fmt.Println(err.Error())

			c.JSON(http.StatusInternalServerError, SERVER_RE{Status: "error", Reply: "failed to process description"})

			return

		}

		c.JSON(http.StatusOK, SERVER_RE{Status: "success", Reply: Encode(answer)})
	})

	return router, nil

}

*/
