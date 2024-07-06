package dbquery

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

var mediaPath string = "./data/media/"

var videoPath string = "./data/media_video/"

func UploadVideo(c *gin.Context, file *multipart.FileHeader, filename string, new_filename string) error {

	/*

		TODO:
			save media index

	*/

	upload_path := path.Join(videoPath, new_filename)

	err := c.SaveUploadedFile(file, upload_path)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	return nil
}

func DownloadVideo(c *gin.Context, watchId string) error {

	download_path := path.Join(videoPath, watchId)

	if _, err := os.Stat(download_path); err != nil {

		return err

	}

	c.Header("Content-Type", "video/mp4")

	c.File(download_path)

	return nil
}
