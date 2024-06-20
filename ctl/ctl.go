package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	pkgauth "github.com/seantywork/sorrylinus-again/pkg/auth"
	pkgstream "github.com/seantywork/sorrylinus-again/pkg/stream"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
)

func CreateServer() *gin.Engine {

	genserver := gin.Default()

	store := sessions.NewCookieStore([]byte("SOLIAGAIN"))

	genserver.Use(sessions.Sessions("SOLIAGAIN", store))

	ConfigureRuntime(genserver)

	RegisterRoutes(genserver)

	return genserver

}

func ConfigureRuntime(e *gin.Engine) {

	e.MaxMultipartMemory = CONF.MaxFileSize

	pkgauth.DEBUG = CONF.Debug

	pkgstream.EXTERNAL_URL = CONF.ExternalUrl

	for i := 0; i < len(CONF.Stream.TurnServerAddr); i++ {

		tmp := struct {
			Addr string `json:"addr"`
			Id   string `json:"id"`
			Pw   string `json:"pw"`
		}{
			Addr: CONF.Stream.TurnServerAddr[i].Addr,
			Id:   CONF.Stream.TurnServerAddr[i].Id,
			Pw:   CONF.Stream.TurnServerAddr[i].Pw,
		}

		pkgstream.TURN_SERVER_ADDR = append(pkgstream.TURN_SERVER_ADDR, tmp)
	}

	pkgstream.RTCP_PLI_INTERVAL = time.Second * time.Duration(CONF.Stream.RtcpPLIInterval)
	pkgstream.UPLOAD_DEST = CONF.Stream.UploadDest
	pkgstream.EXTENSION_ALLOWLIST = CONF.Stream.ExtAllowList

	pkgstream.UDP_BUFFER_BYTE_SIZE = CONF.Stream.UdpBufferByteSize

	pkgstream.SIGNAL_ADDR = CONF.ServeAddr
	pkgstream.SIGNAL_PORT = fmt.Sprintf("%d", CONF.Stream.SignalPort)
	pkgstream.SIGNAL_PORT_EXTERNAL = fmt.Sprintf("%d", CONF.Stream.SignalPortExternal)

	pkgstream.RTP_RECEIVE_ADDR = CONF.ServeAddr
	pkgstream.RTP_RECEIVE_PORT = fmt.Sprintf("%d", CONF.Stream.RtpReceivePort)
	pkgstream.RTP_RECEIVE_PORT_EXTERNAL = fmt.Sprintf("%d", CONF.Stream.RtpReceivePortExternal)

	pkgutils.USE_COMPRESS = CONF.Utils.UseCompress

}

func RegisterRoutes(e *gin.Engine) {

	// base

	e.LoadHTMLGlob("view/*")

	e.Static("/public", "./public")

	e.GET("/", GetIndex)

	e.GET("/signin", GetSigninIndex)

	e.GET("/api/oauth2/google/signin", pkgauth.OauthGoogleLogin)

	e.GET("/oauth2/google/callback", pkgauth.OauthGoogleCallback)

	pkgauth.InitAuth()

	e.GET("/api/turn/address", pkgstream.GetTurnServeAddr)

	// stream

	// cctv

	e.GET("/cctv", pkgstream.GetCCTVIndex)

	e.POST("/api/cctv/create", pkgstream.PostCCTVCreate)

	e.POST("/api/cctv/delete", pkgstream.PostCCTVDelete)

	go pkgstream.InitRTMPServer()

	// video

	e.GET("/video", pkgstream.GetVideoIndex)

	e.GET("/api/video/watch", pkgstream.GetVideoWatchPage)

	e.POST("/api/video/upload", pkgstream.PostVideoUpload)

	e.GET("/api/video/watch/c/:contentId", pkgstream.GetVideoWatchContentByID)

	// peers

	e.GET("/peers", pkgstream.GetPeersIndex)

	e.GET("/api/peers/signal/address", pkgstream.GetPeersSignalAddress)

	go pkgstream.InitPeersSignalOn("/ch/peers/signal")

}
