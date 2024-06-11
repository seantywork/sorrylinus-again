package stream

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
)

func createPeersSignalHandlerForWS() {

	http.HandleFunc("/signal", peerSignalHandler)

	log.Fatal(http.ListenAndServe(*ADDR, nil))

}

func createRoomSignalHandlerForWS() {

	http.HandleFunc("/signal", roomSignalHanler)

	log.Fatal(http.ListenAndServe(*ADDR, nil))

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

	uid := sinfo.UserID

	old_c, okay := USER_SIGNAL[uid]

	if okay {

		log.Printf("uid: %s, exists, removing previous conn\n", uid)

		old_c.Close()

		delete(USER_SIGNAL, uid)

	}

	for k, v := range USER_SIGNAL {

		new_uinfo := SIGNAL_INFO{
			Command: "ADDUSER",
			UserID:  uid,
		}

		v.WriteJSON(&new_uinfo)

		old_uinfo := SIGNAL_INFO{

			Command: "ADDUSER",
			UserID:  k,
		}

		c.WriteJSON(&old_uinfo)

	}

	USER_SIGNAL[uid] = c

	for {

		geninfo := SIGNAL_INFO{}

		c.ReadJSON(&geninfo)

	}

}

func roomSignalHanler(w http.ResponseWriter, r *http.Request) {
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
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Print(err)
		return
	}

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
		if i == nil {
			return
		}

		candidateStringEnc := Encode(i.ToJSON())

		if writeErr := c.WriteJSON(&SIGNAL_INFO{
			Command: "candidate",
			Data:    candidateStringEnc,
		}); writeErr != nil {
			log.Println(writeErr)
		}
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

		log.Println(message.Command)

		switch message.Command {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}

			Decode(message.Data, &candidate)

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

			Decode(message.Data, &answer)

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
