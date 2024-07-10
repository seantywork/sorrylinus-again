package stream

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pkgdbq "github.com/seantywork/sorrylinus-again/pkg/dbquery"
	pkgutils "github.com/seantywork/sorrylinus-again/pkg/utils"
)

var EXTENSION_ALLOWLIST []string

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

	err := pkgdbq.UploadVideo(c, file, f_name, file_name+"."+extension)

	if err != nil {

		fmt.Println(err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RE{Status: "error", Reply: "failed to save"})

		return

	}

	c.JSON(http.StatusOK, SERVER_RE{Status: "success", Reply: "uploaded"})

}

func GetVideoContentByID(c *gin.Context) {

	watchId := c.Param("contentId")

	err := pkgdbq.DownloadVideo(c, watchId)

	if err != nil {

		fmt.Println("path doesn't exist")

		c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid format"})

		return

	}

	fmt.Println("download success")

}
