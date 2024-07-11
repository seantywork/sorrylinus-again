package controller

import "github.com/gin-gonic/gin"

func GetIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

func GetViewSignin(c *gin.Context) {

	c.HTML(200, "signin.html", gin.H{})

}

func GetViewMypage(c *gin.Context) {

	c.HTML(200, "mypage/index.html", gin.H{})

}

func GetViewMypageArticle(c *gin.Context) {

	c.HTML(200, "mypage/article.html", gin.H{})

}

func GetViewMypageVideo(c *gin.Context) {

	c.HTML(200, "mypage/video.html", gin.H{})

}

func GetViewMypageRoom(c *gin.Context) {

	c.HTML(200, "mypage/room.html", gin.H{})

}

func GetViewContentArticle(c *gin.Context) {

	c.HTML(200, "content/article.html", gin.H{})

}

func GetViewContentVideo(c *gin.Context) {

	c.HTML(200, "content/video.html", gin.H{})

}

func GetViewRoom(c *gin.Context) {

	c.HTML(200, "room/index.html", gin.H{})

}
