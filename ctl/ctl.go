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

	pkgstream.INTERNAL_URL = CONF.InternalUrl

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

	pkgstream.PEERS_SIGNAL_PATH = CONF.Stream.PeerSignalAddr

	pkgstream.RTCP_PLI_INTERVAL = time.Second * time.Duration(CONF.Stream.RtcpPLIInterval)
	pkgstream.EXTENSION_ALLOWLIST = CONF.Stream.ExtAllowList

	pkgstream.UDP_BUFFER_BYTE_SIZE = CONF.Stream.UdpBufferByteSize
	pkgstream.UDP_MUX_PORT = CONF.Stream.UdpMuxPort
	pkgstream.UDP_EPHEMERAL_PORT_MIN = CONF.Stream.UdpEphemeralPortMin
	pkgstream.UDP_EPHEMERAL_PORT_MAX = CONF.Stream.UdpEphemeralPortMax

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

	e.GET("/signin", GetViewSignin)

	e.GET("/mypage", GetViewMypage)

	e.GET("/content/article/:articleId", GetViewContentArticle)

	e.GET("/content/peers/:peersId", GetViewContentPeers)

	e.GET("/content/video/:videoId", GetViewContentVideo)

	// auth

	e.GET("/api/oauth2/google/signin", pkgauth.OauthGoogleLogin)

	e.GET("/oauth2/google/callback", pkgauth.OauthGoogleCallback)

	// e.POST("/api/auth/user/add", pkgauth.UserAdd)

	// e.GET("/api/auth/signin", pkgauth.Login)

	// e.GET("/api/auth/signout", pkgauth.Logout)

	pkgauth.InitAuth()

	// edition

	// e.POST("/api/article/upload", pkgedition.PostArticleUpload)

	// e.GET("/api/article/c/:contentId", pkgeditioni.GetArticleContentById)

	// e.POST("/api/image/upload", pkgedition.PostImageUpload)

	// e.GET("/api/image/c/:contentId", pkgedition.GetImageContentById)

	// stream

	pkgstream.InitWebRTCApi()

	e.POST("/api/cctv/create", pkgstream.PostCCTVCreate)

	e.POST("/api/cctv/delete", pkgstream.PostCCTVDelete)

	go pkgstream.InitRTMPServer()

	e.POST("/api/video/upload", pkgstream.PostVideoUpload)

	e.GET("/api/video/c/:contentId", pkgstream.GetVideoContentByID)

	e.GET("/api/peers/signal/address", pkgstream.GetPeersSignalAddress)

	pkgstream.AddSignalHandler(CONF.Stream.PeerSignalAddr, pkgstream.RoomSignalHandler)

	go pkgstream.StartSignalHandler()
}
