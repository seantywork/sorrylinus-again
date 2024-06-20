package stream

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
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

func runPeersSignalHandlerForWS(peerSignalPath string) {

	go func() {

		for range time.NewTicker(time.Second * 3).C {
			dispatchKeyFrame()
		}

	}()

	http.HandleFunc(peerSignalPath, roomSignalHandler)

	signal_addr := SIGNAL_ADDR + ":" + SIGNAL_PORT

	log.Fatal(http.ListenAndServe(signal_addr, nil))

}

func signalPeerConnections() {
	listLock.Lock()

	defer func() {
		listLock.Unlock()
		dispatchKeyFrame()
	}()

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				signalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}

func attemptSync() bool {

	for i := range peerConnections {
		if peerConnections[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
			peerConnections = append(peerConnections[:i], peerConnections[i+1:]...)
			return true // We modified the slice, start from the beginning
		}

		// map of sender we already are seanding, so we don't double send
		existingSenders := map[string]bool{}

		for _, sender := range peerConnections[i].peerConnection.GetSenders() {
			if sender.Track() == nil {
				continue
			}

			existingSenders[sender.Track().ID()] = true

			// If we have a RTPSender that doesn't map to a existing track remove and signal
			if _, ok := trackLocals[sender.Track().ID()]; !ok {
				if err := peerConnections[i].peerConnection.RemoveTrack(sender); err != nil {
					return true
				}
			}
		}

		// Don't receive videos we are sending, make sure we don't have loopback
		for _, receiver := range peerConnections[i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			existingSenders[receiver.Track().ID()] = true
		}

		// Add all track we aren't sending yet to the PeerConnection
		for trackID := range trackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				if _, err := peerConnections[i].peerConnection.AddTrack(trackLocals[trackID]); err != nil {
					return true
				}
			}
		}

		offer, err := peerConnections[i].peerConnection.CreateOffer(nil)
		if err != nil {
			return true
		}

		if err = peerConnections[i].peerConnection.SetLocalDescription(offer); err != nil {
			return true
		}

		offerStringEnc := pkgutils.Encode(offer)

		if err = peerConnections[i].websocket.WriteJSON(&SIGNAL_INFO{
			Command: "offer",
			Data:    offerStringEnc,
		}); err != nil {
			return true
		}
	}

	return false
}

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

func peerSignalHandler(w http.ResponseWriter, r *http.Request) {

	UPGRADER.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := UPGRADER.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgraded")
		return
	}

	c.SetReadDeadline(time.Time{})

	var sinfo SIGNAL_INFO

	err = c.ReadJSON(&sinfo)

	if err != nil {

		log.Printf("failed to read uinfo: %s\n", err.Error())

		return

	}

	uid := sinfo.Data

	old_c, okay := USER_SIGNAL[uid]

	if okay {

		log.Printf("uid: %s, exists, removing previous conn\n", uid)

		old_c.Close()

		delete(USER_SIGNAL, uid)

	}

	for k, v := range USER_SIGNAL {

		new_uinfo := SIGNAL_INFO{
			Command: "ADDUSER",
			Data:    uid,
		}

		v.WriteJSON(&new_uinfo)

		old_uinfo := SIGNAL_INFO{

			Command: "ADDUSER",
			Data:    k,
		}

		c.WriteJSON(&old_uinfo)

	}

	USER_SIGNAL[uid] = c

	for {

		geninfo := SIGNAL_INFO{}

		c.ReadJSON(&geninfo)

	}

}

func roomSignalHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request to Websocket

	UPGRADER.CheckOrigin = func(r *http.Request) bool { return true }

	unsafeConn, err := UPGRADER.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	c := &threadSafeWriter{unsafeConn, sync.Mutex{}}

	// When this frame returns close the Websocket
	defer c.Close() //nolint

	// Create new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:       []string{TURN_SERVER_ADDR[0].Addr},
				Username:   TURN_SERVER_ADDR[0].Id,
				Credential: TURN_SERVER_ADDR[0].Pw,
			},
		},
	})

	if err != nil {
		log.Print(err)
		return
	}

	log.Print("new peerconnection added")

	// When this frame returns close the PeerConnection
	defer peerConnection.Close() //nolint

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Print(err)
			return
		}
	}

	// Add our new PeerConnection to global list
	listLock.Lock()
	peerConnections = append(peerConnections, peerConnectionState{peerConnection, c})
	listLock.Unlock()

	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {

		log.Printf("got candidate\n")

		if i == nil {
			return
		}

		candidateStringEnc := pkgutils.Encode(i.ToJSON())

		if writeErr := c.WriteJSON(&SIGNAL_INFO{
			Command: "candidate",
			Data:    candidateStringEnc,
		}); writeErr != nil {
			log.Println(writeErr)
		}

		log.Printf("sent candidate\n")
	})

	// If PeerConnection is closed remove it from global list
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			signalPeerConnections()
		default:
			log.Printf("on connection state change: %s \n", p.String())
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := addTrack(t)
		defer removeTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	// Signal for the new PeerConnection
	signalPeerConnections()

	message := &SIGNAL_INFO{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}

		log.Printf("got message: %s\n", message.Command)

		switch message.Command {
		case "candidate":

			candidate := webrtc.ICECandidateInit{}

			pkgutils.Decode(message.Data, &candidate)

			/*
				if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
					log.Println(err)
					return
				}
			*/

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Println(err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}

			pkgutils.Decode(message.Data, &answer)

			/*
				if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
					log.Println(err)
					return
				}
			*/

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
