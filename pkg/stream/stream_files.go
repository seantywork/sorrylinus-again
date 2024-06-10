package stream

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	pkgutils "github.com/lineworld-lab/go-tv/pkg/utils"
)

var MAX_FILE_SIZE int64 = 838860800 // 800Mib

var UPLOAD_DEST string = "./upload/"

var EXTENSION_ALLOWLIST []string = []string{
	"mp4",
}

func CreateStreamServerForFiles() (*gin.Engine, error) {

	router := CreateGenericServer()

	router.MaxMultipartMemory = MAX_FILE_SIZE

	router.GET("/", func(c *gin.Context) {

		c.HTML(200, "files.html", gin.H{
			"title": "Files",
		})

	})

	router.GET("/watch", func(c *gin.Context) {

		c.HTML(200, "files_watch.html", gin.H{
			"title": "Files Watch",
		})

	})

	router.POST("/files/upload", func(c *gin.Context) {

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

	})

	router.GET("/files/d/:watchId", func(c *gin.Context) {

		watchId := c.Param("watchId")

		extension := "mp4"

		upload_path := path.Join(UPLOAD_DEST, watchId+"."+extension)

		if _, err := os.Stat(upload_path); err != nil {

			fmt.Println("path doesn't exist")

			c.JSON(http.StatusBadRequest, SERVER_RE{Status: "error", Reply: "invalid format"})

			return

		}

		c.Header("Content-Type", "video/mp4")

		c.File(upload_path)

	})

	return router, nil
}
