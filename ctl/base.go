package controller

import "github.com/gin-gonic/gin"

func GetIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

func GetViewSignin(c *gin.Context) {

	c.HTML(200, "signin.html", gin.H{})

}

func GetViewMypage(c *gin.Context) {

	c.HTML(200, "mypage.html", gin.H{})

}

func GetViewContentArticle(c *gin.Context) {

	c.HTML(200, "content/article.html", gin.H{})

}

func GetViewContentPeers(c *gin.Context) {

	c.HTML(200, "content/peers.html", gin.H{})

}

func GetViewContentVideo(c *gin.Context) {

	c.HTML(200, "content/video.html", gin.H{})

}
