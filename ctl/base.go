package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pkgauth "github.com/seantywork/sorrylinus-again/pkg/auth"
	"github.com/seantywork/sorrylinus-again/pkg/com"
	"github.com/seantywork/sorrylinus-again/pkg/dbquery"
)

type EntryStruct struct {
	Entry []struct {
		Title string
		Id    string
	} `json:"entry"`
}

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

func GetViewContentArticle(c *gin.Context) {

	c.HTML(200, "content/article.html", gin.H{})

}

func GetViewContentVideo(c *gin.Context) {

	c.HTML(200, "content/video.html", gin.H{})

}

func GetViewRoom(c *gin.Context) {

	c.HTML(200, "room/index.html", gin.H{})

}

func GetMediaEntry(c *gin.Context) {

	entry := EntryStruct{}

	em, err := dbquery.GetEntryForMedia()

	if err != nil {

		fmt.Printf("get content entry: failed to retrieve: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, com.SERVER_RE{Status: "error", Reply: "failed to retrieve content entry"})

		return

	}

	for k, v := range em {

		if v.Type == "article" {

			entry.Entry = append(entry.Entry, struct {
				Title string
				Id    string
			}{

				Title: v.PlainName,
				Id:    k,
			})

		} else if v.Type == "video" {

			entry.Entry = append(entry.Entry, struct {
				Title string
				Id    string
			}{

				Title: v.PlainName + "." + v.Extension,
				Id:    k,
			})

		} else {

			continue
		}

	}

	jb, err := json.Marshal(entry)

	if err != nil {

		fmt.Printf("get content entry: failed to marshal: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, com.SERVER_RE{Status: "error", Reply: "failed to retrieve content entry"})

		return

	}

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: string(jb)})

}
