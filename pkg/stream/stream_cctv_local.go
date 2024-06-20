package stream

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	webrtcv4 "github.com/pion/webrtc/v4"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
)

var UDP_BUFFER_BYTE_SIZE int

var RECV_STARTED int = 0

func GetCCTVLocalIndex(c *gin.Context) {

	c.HTML(200, "cctv_local.html", gin.H{
		"title": "CCTV",
	})
}

func PostCCTVLocalOffer(c *gin.Context) {

	if RECV_STARTED == 1 {

		fmt.Println("recv already started")

		c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid request"})

		return
	}

	var offerjson CLIENT_REQ

	if err := c.BindJSON(&offerjson); err != nil {

		fmt.Println("failed to get request body")

		c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	offer_out := make(chan string)

	go startCCTVReceiver(offerjson.Data, offer_out)

	offer_out_str := <-offer_out

	c.JSON(http.StatusOK, SERVER_RE{Status: "success", Reply: offer_out_str})

}
func startCCTVReceiver(offer_in string, offer_out chan string) {

	peerConnection, err := webrtcv4.NewPeerConnection(webrtcv4.Configuration{
		ICEServers: []webrtcv4.ICEServer{
			{
				URLs: []string{TURN_SERVER_ADDR},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	var udp_port int

	fmt.Sscanf(RTP_RECEIVE_PORT, "%d", &udp_port)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(RTP_RECEIVE_ADDR), Port: udp_port})
	if err != nil {
		panic(err)
	}

	// Increase the UDP receive buffer size
	// Default UDP buffer sizes vary on different operating systems
	//bufferSize := UDP_BUFFER_BYTE_SIZE // 300KB
	//err = listener.SetReadBuffer(bufferSize)
	//if err != nil {
	//	panic(err)
	//}

	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a video track
	videoTrack, err := webrtcv4.NewTrackLocalStaticRTP(webrtcv4.RTPCodecCapability{MimeType: webrtcv4.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtcv4.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtcv4.ICEConnectionStateFailed {
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
		}
	})

	// Wait for the offer to be pasted
	offer := webrtcv4.SessionDescription{}
	pkgutils.Decode(offer_in, &offer)

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtcv4.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser

	localdesc := pkgutils.Encode(peerConnection.LocalDescription())

	offer_out <- localdesc

	// Read RTP packets forever and send them to the WebRTC Client
	inboundRTPPacket := make([]byte, 1600) // UDP MTU

	RECV_STARTED = 1

	new_conn, err := listener.Accept()

	if err != nil {

		panic(err)

	}

	for {

		n, err := new_conn.Read(inboundRTPPacket)
		//n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {

			RECV_STARTED = 0

			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// The peerConnection has been closed.

				RECV_STARTED = 0
				return
			}

			RECV_STARTED = 0

			panic(err)
		}
	}
}
