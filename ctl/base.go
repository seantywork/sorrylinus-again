package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pkgauth "github.com/seantywork/sorrylinus-again/pkg/auth"
	"github.com/seantywork/sorrylinus-again/pkg/com"
)

func GetIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

func GetViewSignin(c *gin.Context) {

	c.HTML(200, "signin.html", gin.H{})

}

func GetViewMypage(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("view my page: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	c.HTML(200, "mypage/index.html", gin.H{})

}

func GetViewMypageArticle(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("view my page: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	c.HTML(200, "mypage/article.html", gin.H{})

}

func GetViewMypageVideo(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("view my page: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	c.HTML(200, "mypage/video.html", gin.H{})

}

func GetViewMypageRoom(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("view my page: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	c.HTML(200, "mypage/room.html", gin.H{})

}

func GetBase(c *gin.Context) {

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
