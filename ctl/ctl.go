package controller

import (
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

	pkgstream.TURN_SERVER_ADDR = CONF.Stream.TurnServerAddr
	pkgstream.RTCP_PLI_INTERVAL = time.Second * time.Duration(CONF.Stream.RtcpPLIInterval)
	pkgstream.UPLOAD_DEST = CONF.Stream.UploadDest
	pkgstream.EXTENSION_ALLOWLIST = CONF.Stream.ExtAllowList

	pkgstream.UDP_BUFFER_BYTE_SIZE = CONF.Stream.UdpBufferByteSize

	pkgstream.SIGNAL_ADDR = CONF.ServeAddr
	pkgstream.SIGNAL_PORT = CONF.Stream.SignalPort

	pkgstream.RTP_RECEIVE_ADDR = CONF.ServeAddr
	pkgstream.RTP_RECEIVE_PORT = CONF.Stream.RtpReceivePort

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
