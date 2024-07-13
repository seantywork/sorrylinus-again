package com

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type CLIENT_REQ struct {
	Data string `json:"data"`
}

type SERVER_RE struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

type RT_REQ_DATA struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type RT_RESP_DATA struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

var CHANNEL_ADDR string

var CHANNEL_PORT string

var CHANNEL_PORT_EXTERNAL string

var USER_CHANNEL = make(map[string]*websocket.Conn)

var CH_CALLBACKS []func()

var ListLock sync.RWMutex

type ThreadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

func AddChannelHandler(channelPath string, channelHandler func(w http.ResponseWriter, r *http.Request)) {

	http.HandleFunc(channelPath, channelHandler)

}

func AddChannelCallback(channelFunction func()) {

	CH_CALLBACKS = append(CH_CALLBACKS, channelFunction)
}

func StartAllChannelHandlers() {

	callback_count := len(CH_CALLBACKS)

	for i := 0; i < callback_count; i++ {

		go CH_CALLBACKS[i]()

	}

	channel_addr := CHANNEL_ADDR + ":" + CHANNEL_PORT

	log.Fatal(http.ListenAndServe(channel_addr, nil))

}
