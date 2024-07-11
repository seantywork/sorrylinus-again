package edition

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pkgauth "github.com/seantywork/sorrylinus-again/pkg/auth"
	"github.com/seantywork/sorrylinus-again/pkg/com"
	"github.com/seantywork/sorrylinus-again/pkg/dbquery"
	_ "github.com/seantywork/sorrylinus-again/pkg/dbquery"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
)

type ArticleInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func PostArticleUpload(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("article upload: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	var req com.CLIENT_REQ

	var a_info ArticleInfo

	if err := c.BindJSON(&req); err != nil {

		fmt.Printf("article upload: failed to bind: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return
	}

	err := json.Unmarshal([]byte(req.Data), &a_info)

	if err != nil {

		fmt.Printf("article upload: failed to unmarshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	new_file_name, _ := pkgutils.GetRandomHex(32)

	plain_name := pkgauth.SanitizePlainNameValue(a_info.Title)

	err = dbquery.UploadArticle(a_info.Content, plain_name, new_file_name)

	if err != nil {

		fmt.Printf("article upload: failed to upload: %s", err.Error())

		c.JSON(http.StatusInternalServerError, com.SERVER_RE{Status: "error", Reply: "failed to upload"})

		return
	}

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: "uploaded"})

}

func PostArticleDelete(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("article delete: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	fmt.Println("delete article")

	var req com.CLIENT_REQ

	if err := c.BindJSON(&req); err != nil {

		fmt.Printf("article delete: failed to bind: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return
	}

	if !pkgauth.VerifyCodeNameValue(req.Data) {

		fmt.Printf("article name verification failed: %s\n", req.Data)

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	err := dbquery.DeleteArticle(req.Data)

	if err != nil {

		fmt.Printf("article delete: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, com.SERVER_RE{Status: "error", Reply: "failed delete"})

		return

	}

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: "deleted"})

}

func GetArticleContentById(c *gin.Context) {

	watchId := c.Param("contentId")

	if !pkgauth.VerifyCodeNameValue(watchId) {

		fmt.Printf("get article: illegal: %s\n", watchId)

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	content, err := dbquery.GetArticle(watchId)

	if err != nil {

		fmt.Printf("failed to get article: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return
	}

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: content})

}

func PostImageUpload(c *gin.Context) {

	_, my_type, _ := pkgauth.WhoAmI(c)

	if my_type != "admin" {

		fmt.Printf("image upload: not admin\n")

		c.JSON(http.StatusForbidden, com.SERVER_RE{Status: "error", Reply: "you're not admin"})

		return

	}

	file, _ := c.FormFile("file")

	f_name := file.Filename

	f_name_list := strings.Split(f_name, ".")

	f_name_len := len(f_name_list)

	if f_name_len < 1 {

		fmt.Println("no extension specified")

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return
	}

	v_fname := pkgauth.SanitizePlainNameValue(f_name_list[0])

	extension := f_name_list[f_name_len-1]

	if !pkgutils.CheckIfSliceContains[string](EXTENSION_ALLOWLIST, extension) {

		fmt.Println("extension not allowed")

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	fmt.Printf("received: %s, size: %d\n", file.Filename, file.Size)

	file_name, _ := pkgutils.GetRandomHex(32)

	err := dbquery.UploadImage(c, file, v_fname, file_name, extension)

	if err != nil {

		fmt.Println(err.Error())

		c.JSON(http.StatusInternalServerError, com.SERVER_RE{Status: "error", Reply: "failed to save"})

		return

	}

	c.JSON(http.StatusOK, com.SERVER_RE{Status: "success", Reply: file_name})

}

func GetImageContentById(c *gin.Context) {

	watchId := c.Param("contentId")

	if !pkgauth.VerifyCodeNameValue(watchId) {

		fmt.Printf("download image: illegal: %s\n", watchId)

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	err := dbquery.DownloadImage(c, watchId)

	if err != nil {

		fmt.Printf("download image: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	fmt.Println("image download success")
}
