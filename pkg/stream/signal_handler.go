package stream

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var USER_SIGNAL = make(map[string]*websocket.Conn)

var ADDR = flag.String("addr", "0.0.0.0:8082", "service address")

var UPGRADER = websocket.Upgrader{}

type SIGNAL_INFO struct {
	Command string `json:"command"`
	UserID  string `json:"user_id"`
}

func createSignalHandlerForWS() {

	http.HandleFunc("/signal", signalHandler)

	log.Fatal(http.ListenAndServe(*ADDR, nil))

}

func signalHandler(w http.ResponseWriter, r *http.Request) {

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
