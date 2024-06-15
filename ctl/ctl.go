package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	pkgstream "github.com/seantywork/sorrylinus-again/pkg/stream"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
)

func CreateServer() *gin.Engine {

	genserver := gin.Default()

	genserver.MaxMultipartMemory = CONF.MaxFileSize

	pkgstream.EXTERNAL_URL = CONF.ExternalUrl

	pkgstream.TURN_SERVER_ADDR = CONF.Stream.TurnServerAddr
	pkgstream.RTCP_PLI_INTERVAL = time.Second * CONF.Stream.RtcpPLIInterval
	pkgstream.UPLOAD_DEST = CONF.Stream.UploadDest
	pkgstream.EXTENSION_ALLOWLIST = CONF.Stream.ExtAllowList

	pkgstream.SIGNAL_ADDR = CONF.ServeAddr
	pkgstream.SIGNAL_PORT = CONF.Stream.SignalPort

	pkgstream.RTP_RECEIVE_ADDR = CONF.ServeAddr
	pkgstream.RTP_RECEIVE_PORT = CONF.Stream.RtpReceivePort

	pkgutils.USE_COMPRESS = CONF.Utils.UseCompress

	// base

	genserver.LoadHTMLGlob("view/*")

	genserver.Static("/public", "./public")

	// stream

	// cctv

	genserver.GET("/cctv", pkgstream.GetCCTVIndex)

	genserver.GET("/cctv/turn/address", pkgstream.GetCCTVTurnServeAddr)

	genserver.POST("/cctv/offer", pkgstream.PostCCTVOffer)

	// video

	genserver.GET("/video", pkgstream.GetVideoIndex)

	genserver.GET("/video/watch", pkgstream.GetVideoWatchPage)

	genserver.POST("/video/upload", pkgstream.PostVideoUpload)

	genserver.GET("/video/watch/c/:contentId", pkgstream.GetVideoWatchContentByID)

	// peers

	genserver.GET("/peers", pkgstream.GetPeersIndex)

	genserver.GET("/peers/signal/address", pkgstream.GetPeersSignalAddress)

	go pkgstream.InitPeersSignalOn("/peers/signal")

	// utils

	return genserver

}
