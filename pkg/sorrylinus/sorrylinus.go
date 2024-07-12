package sorrylinus

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/seantywork/sorrylinus-again/pkg/com"
)

var SOLI_FRONT_ADDR string

var SOLI_SIGNAL_PATH string

var EXTERNAL_URL string

var TIMEOUT_SEC int

type RT_REQ_DATA struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type RT_RESP_DATA struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

var UREG map[*websocket.Conn]string

var UPGRADER = websocket.Upgrader{}

func GetSoliSignalAddress(c *gin.Context) {

	s_addr := EXTERNAL_URL + ":" + com.CHANNEL_PORT_EXTERNAL + SOLI_SIGNAL_PATH

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: s_addr})

}

func SoliSignalHandler(w http.ResponseWriter, r *http.Request) {

	c, _, err := websocket.DefaultDialer.Dial(SOLI_FRONT_ADDR, nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	timeout_iter_count := 0

	timeout_iter := TIMEOUT_SEC * 10

	ticker := time.NewTicker(100 * time.Millisecond)

	received_auth := make(chan RT_REQ_DATA)

	got_auth := 0

	c_ret := RT_RESP_DATA{}

	var req RT_REQ_DATA

	go func() {

		auth_req := RT_REQ_DATA{}

		err := c.ReadJSON(&auth_req)

		if err != nil {

			log.Fatal("read auth:", err)
			return
		}

		received_auth <- auth_req

	}()

	for got_auth == 0 {

		select {

		case <-ticker.C:

			if timeout_iter_count <= timeout_iter {

				timeout_iter_count += 1

			} else {

				log.Fatal("read auth timed out:", err)

				return
			}

		case a := <-received_auth:

			req = a

			got_auth = 1

			break
		}

	}

	resp, err := RoundTrip(&req)

	if err != nil {

		c_ret.Status = "error"

		c.WriteJSON(c_ret)

		log.Fatal("rt: auth error:", err)

		return
	}

	if resp.Status != "success" {

		c_ret.Status = "error"

		c_ret.Data = "authentication failed"

		c.WriteJSON(c_ret)

		return

	}

}
