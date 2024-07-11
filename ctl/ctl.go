package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	pkgauth "github.com/seantywork/sorrylinus-again/pkg/auth"
	pkgcom "github.com/seantywork/sorrylinus-again/pkg/com"
	pkgedition "github.com/seantywork/sorrylinus-again/pkg/edition"
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

	pkgcom.CHANNEL_ADDR = CONF.ServeAddr
	pkgcom.CHANNEL_PORT = fmt.Sprintf("%d", CONF.Com.ChannelPort)
	pkgcom.CHANNEL_PORT_EXTERNAL = fmt.Sprintf("%d", CONF.Com.ChannelPortExternal)

	pkgedition.EXTENSION_ALLOWLIST = CONF.Edition.ExtAllowList

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

	pkgstream.UDP_BUFFER_BYTE_SIZE = CONF.Stream.UdpBufferByteSize
	pkgstream.UDP_MUX_PORT = CONF.Stream.UdpMuxPort
	pkgstream.UDP_EPHEMERAL_PORT_MIN = CONF.Stream.UdpEphemeralPortMin
	pkgstream.UDP_EPHEMERAL_PORT_MAX = CONF.Stream.UdpEphemeralPortMax

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

	e.GET("/mypage/article", GetViewMypageArticle)

	e.GET("/mypage/video", GetViewMypageVideo)

	e.GET("/mypage/room", GetViewMypageRoom)

	e.GET("/api/base", GetBase)

	e.GET("/content/article/:articleId", GetViewContentArticle)

	e.GET("/content/video/:videoId", GetViewContentVideo)

	e.GET("/room/:roomId", GetViewRoom)

	// auth

	e.GET("/api/oauth2/google/signin", pkgauth.OauthGoogleLogin)

	e.GET("/oauth2/google/callback", pkgauth.OauthGoogleCallback)

	e.POST("/api/auth/user/add", pkgauth.UserAdd)

	e.POST("/api/auth/user/remove", pkgauth.UserRemove)

	e.GET("/api/auth/signin", pkgauth.Login)

	e.GET("/api/auth/signout", pkgauth.Logout)

	pkgauth.InitAuth()

	// sorrylinus

	// e.POST("/api/sorrylinus/connect", pkgsoli.Connect)
	// e.POST("/api/sorrylinus/disconnect", pkgsoli.Disconnect)
	// e.POST("/api/sorrylinus/rt", pkgsoli.RoundTrip)

	// edition

	e.POST("/api/article/upload", pkgedition.PostArticleUpload)

	e.POST("/api/article/delete", pkgedition.PostArticleDelete)

	e.GET("/api/article/c/:contentId", pkgedition.GetArticleContentById)

	e.POST("/api/image/upload", pkgedition.PostImageUpload)

	e.GET("/api/image/c/:contentId", pkgedition.GetImageContentById)

	e.POST("/api/video/upload", pkgedition.PostVideoUpload)

	e.POST("/api/video/delete", pkgedition.PostVideoDelete)

	e.GET("/api/video/c/:contentId", pkgedition.GetVideoContentByID)

	// stream

	pkgstream.InitWebRTCApi()

	e.POST("/api/cctv/open", pkgstream.PostCCTVOpen)

	e.POST("/api/cctv/close", pkgstream.PostCCTVClose)

	go pkgstream.InitRTMPServer()

	e.GET("/api/peers/signal/address", pkgstream.GetPeersSignalAddress)

	// channel

	pkgcom.AddChannelHandler(CONF.Stream.PeerSignalAddr, pkgstream.RoomSignalHandler)

	pkgcom.AddChannelCallback(pkgstream.SignalDispatcher)

	go pkgcom.StartAllChannelHandlers()
}
