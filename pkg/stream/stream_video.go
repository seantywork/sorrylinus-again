package stream

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
)

var UPLOAD_DEST string = "./upload/"

var EXTENSION_ALLOWLIST []string = []string{
	"mp4",
}

func GetVideoIndex(c *gin.Context) {

	c.HTML(200, "video.html", gin.H{
		"title": "Video",
	})

}

func GetVideoWatchPage(c *gin.Context) {

	c.HTML(200, "video_watch.html", gin.H{
		"title": "Video Watch",
	})

}

func PostVideoUpload(c *gin.Context) {

	file, _ := c.FormFile("file")

	f_name := file.Filename

	f_name_list := strings.Split(f_name, ".")

	f_name_len := len(f_name_list)

	if f_name_len < 1 {

		fmt.Println("no extension specified")

		c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid format"})

		return
	}

	extension := f_name_list[f_name_len-1]

	if !pkgutils.CheckIfSliceContains[string](f_name_list, extension) {

		fmt.Println("extension not allowed")

		c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	fmt.Printf("received: %s, size: %d\n", file.Filename, file.Size)

	file_name, _ := pkgutils.GetRandomHex(32)

	upload_path := path.Join(UPLOAD_DEST, file_name+"."+extension)

	c.SaveUploadedFile(file, upload_path)

	c.JSON(http.StatusOK, SERVER_RE{Status: "success", Reply: "uploaded"})

}

func GetVideoWatchContentByID(c *gin.Context) {

	watchId := c.Param("contentId")

	extension := "mp4"

	upload_path := path.Join(UPLOAD_DEST, watchId+"."+extension)

	if _, err := os.Stat(upload_path); err != nil {

		fmt.Println("path doesn't exist")

		c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	c.Header("Content-Type", "video/mp4")

	c.File(upload_path)

}
