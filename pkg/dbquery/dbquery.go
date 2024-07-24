package dbquery

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminStruct struct {
}

type UserStruct struct {
	Passphrase      string `json:"passphrase"`
	DurationSeconds int    `json:"duration_seconds"`
}

type SessionStruct struct {
	Type            string `json:"type"`
	Id              string `json:"id"`
	StartTime       string `json:"start_time"`
	DurationSeconds int    `json:"duration_seconds"`
}

type MediaStruct struct {
	ISPublic  bool     `json:"is_public"`
	CreatedAt string   `json:"created_at"`
	Type      string   `json:"type"`
	Extension string   `json:"extension"`
	PlainName string   `json:"plain_name"`
	Author    string   `json:"author"`
	AllowedId []string `json:"allowed_id"`
}

var adminPath string = "./data/admin/"

var userPath string = "./data/user/"

var sessionPath string = "./data/session/"

var mediaPath string = "./data/media/"

var articlePath string = "./data/media/article/"

var imagePath string = "./data/media/image/"

var videoPath string = "./data/media/video/"

func GetByIdFromAdmin(email string) (*AdminStruct, error) {

	var as AdminStruct

	this_file_path := adminPath + email + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return nil, fmt.Errorf("failed to get from admin: %s", err.Error())
	}

	err = json.Unmarshal(file_b, &as)

	if err != nil {

		return nil, fmt.Errorf("failed to get from admin: %s", err.Error())
	}

	return &as, nil

}

func GetByIdFromUser(id string) (*UserStruct, error) {

	var us UserStruct

	this_file_path := userPath + id + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return nil, fmt.Errorf("failed to get from user: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &us)

	if err != nil {

		return nil, fmt.Errorf("failed to get from user: %s", err.Error())
	}

	return &us, nil

}

func GetBySessionKeyFromSession(session_key string) (*SessionStruct, error) {

	var ss SessionStruct

	this_file_path := sessionPath + session_key + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return nil, fmt.Errorf("session pool: failed to read file: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &ss)

	if err != nil {

		return nil, fmt.Errorf("session pool: failed to marshal: %s", err.Error())

	}

	t_now := time.Now()

	t, _ := time.Parse("2006-01-02-15-04-05", ss.StartTime)

	diff := t_now.Sub(t)

	if ss.DurationSeconds == 0 || (diff.Seconds() < float64(ss.DurationSeconds)) {

		return &ss, nil

	} else {

		err = RemoveSessionKeyFromSession(session_key)

		if err != nil {

			return nil, fmt.Errorf("session: failed to remove: %s", err.Error())

		}

		return nil, fmt.Errorf("id: %s: expired", ss.Id)

	}

}

func GetByMediaKeyFromMedia(media_key string) (*MediaStruct, error) {

	var ms MediaStruct

	this_file_path := mediaPath + media_key + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return nil, fmt.Errorf("media: failed to read file: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &ms)

	if err != nil {

		return nil, fmt.Errorf("media: failed to marshal: %s", err.Error())
	}

	return &ms, nil

}

func GetByIdFromSession(email string) (string, *SessionStruct, error) {

	files, err := os.ReadDir(sessionPath)

	if err != nil {

		return "", nil, fmt.Errorf("session pool: failed to read dir: %s", err.Error())

	}

	for _, f := range files {

		ss := SessionStruct{}

		f_name := f.Name()

		if !strings.Contains(f_name, ".json") {
			continue
		}

		key_name := strings.ReplaceAll(f_name, ".json", "")

		this_file_path := sessionPath + f_name

		file_b, err := os.ReadFile(this_file_path)

		if err != nil {

			return "", nil, fmt.Errorf("session pool: failed to read file: %s", err.Error())

		}

		err = json.Unmarshal(file_b, &ss)

		if err != nil {

			return "", nil, fmt.Errorf("session pool: failed to marshal: %s", err.Error())
		}

		t_now := time.Now()

		t, _ := time.Parse("2006-01-02-15-04-05", ss.StartTime)

		diff := t_now.Sub(t)

		if ss.Id == email {

			if ss.DurationSeconds == 0 || (diff.Seconds() < float64(ss.DurationSeconds)) {

				return key_name, &ss, nil

			} else {

				err = RemoveSessionKeyFromSession(key_name)

				if err != nil {

					return "", nil, fmt.Errorf("session pool: failed to remove: %s", err.Error())

				}

				return "", nil, fmt.Errorf("id: %s: expired", email)

			}

		}

	}

	return "", nil, fmt.Errorf("id: %s: not found", email)
}

func GetEntryForMedia() (map[string]MediaStruct, error) {

	em := make(map[string]MediaStruct)

	files, err := os.ReadDir(mediaPath)

	if err != nil {

		return nil, fmt.Errorf("media entry: failed to read dir: %s", err.Error())

	}

	for _, f := range files {

		ms := MediaStruct{}

		f_name := f.Name()

		if !strings.Contains(f_name, ".json") {
			continue
		}

		key_name := strings.ReplaceAll(f_name, ".json", "")

		this_file_path := mediaPath + f_name

		file_b, err := os.ReadFile(this_file_path)

		if err != nil {

			return nil, fmt.Errorf("media entry: failed to read file: %s", err.Error())

		}

		err = json.Unmarshal(file_b, &ms)

		if err != nil {

			return nil, fmt.Errorf("media entry: failed to marshal: %s", err.Error())
		}

		em[key_name] = ms

	}

	return em, nil
}

func RemoveSessionKeyFromSession(session_key string) error {

	this_file_path := sessionPath + session_key + ".json"

	err := os.Remove(this_file_path)

	if err != nil {

		return err

	}

	return nil
}

func MakeSessionForAdmin(session_key string, email string) error {

	_, ss, _ := GetByIdFromSession(email)

	if ss != nil {

		return fmt.Errorf("valid session already exists for: %s", email)

	}

	new_ss := SessionStruct{}

	t_now := time.Now()

	t_str := t_now.Format("2006-01-02-15-04-05")

	new_ss.Type = "admin"

	new_ss.Id = email

	new_ss.StartTime = t_str

	new_ss.DurationSeconds = 0

	jb, err := json.Marshal(new_ss)

	if err != nil {

		return fmt.Errorf("failed to marshal new session admin: %s", err.Error())
	}

	this_file_path := sessionPath + session_key + ".json"

	err = os.WriteFile(this_file_path, jb, 0644)

	if err != nil {

		return fmt.Errorf("failed to write new session admin: %s", err.Error())
	}

	return nil
}

func GetAllUsers() (map[string]UserStruct, error) {

	var ret_us = make(map[string]UserStruct)

	files, err := os.ReadDir(userPath)

	if err != nil {

		return nil, fmt.Errorf("all users: failed to read dir: %s", err.Error())

	}

	for _, f := range files {

		us := UserStruct{}

		f_name := f.Name()

		if !strings.Contains(f_name, ".json") {
			continue
		}

		key_name := strings.ReplaceAll(f_name, ".json", "")

		this_file_path := userPath + f_name

		file_b, err := os.ReadFile(this_file_path)

		if err != nil {

			return nil, fmt.Errorf("all users: failed to read file: %s", err.Error())

		}

		err = json.Unmarshal(file_b, &us)

		if err != nil {

			return nil, fmt.Errorf("all users: failed to marshal: %s", err.Error())
		}

		ret_us[key_name] = us

	}

	return ret_us, nil

}

func MakeUser(id string, passphrase string, duration_seconds int) error {

	var us UserStruct

	this_file_path := userPath + id + ".json"

	if _, err := os.Stat(this_file_path); err == nil {

		return fmt.Errorf("id: %s: exists", err.Error())

	}

	us.Passphrase = passphrase

	us.DurationSeconds = duration_seconds

	jb, err := json.Marshal(us)

	if err != nil {

		return fmt.Errorf("failed to make user: %s", err.Error())

	}

	err = os.WriteFile(this_file_path, jb, 0644)

	if err != nil {

		return fmt.Errorf("failed to make user: %s", err.Error())
	}

	return nil
}

func RemoveUser(id string) error {

	this_file_path := userPath + id + ".json"

	key, ss, _ := GetByIdFromSession(id)

	if ss == nil {

		fmt.Println("remove user: existing session removed")

		_ = RemoveSessionKeyFromSession(key)

	}

	err := os.Remove(this_file_path)

	if err != nil {

		return err

	}

	return nil
}

func MakeSessionForUser(session_key string, id string, duration_seconds int) error {

	var ss SessionStruct

	this_file_path := sessionPath + session_key + ".json"

	t_now := time.Now()

	t_str := t_now.Format("2006-01-02-15-04-05")

	ss.Type = "user"

	ss.Id = id

	ss.StartTime = t_str

	ss.DurationSeconds = duration_seconds

	jb, err := json.Marshal(ss)

	if err != nil {

		return fmt.Errorf("failed to make session for user: %s", err.Error())

	}

	err = os.WriteFile(this_file_path, jb, 0644)

	if err != nil {

		return fmt.Errorf("failed to make session: %s", err.Error())

	}

	return nil
}

func UploadArticle(author string, content string, plain_name string, new_name string) error {

	ms := MediaStruct{}

	c_time := time.Now()

	c_time_fmt := c_time.Format("2006-01-02-15-04-05")

	ms.ISPublic = true
	ms.Type = "article"
	ms.PlainName = plain_name
	ms.Extension = "json"
	ms.CreatedAt = c_time_fmt
	ms.Author = author

	this_file_path := mediaPath + new_name + ".json"

	this_article_path := articlePath + new_name + ".json"

	content_b := []byte(content)

	jb, err := json.Marshal(ms)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	err = os.WriteFile(this_file_path, jb, 0644)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	err = os.WriteFile(this_article_path, content_b, 0644)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	return nil
}

func DeleteArticle(media_key string) error {

	var ms MediaStruct

	this_file_path := mediaPath + media_key + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return fmt.Errorf("failed to delete article: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &ms)

	if err != nil {

		return fmt.Errorf("failed to delete article: marshal: %s", err.Error())
	}

	this_article_path := articlePath + media_key + "." + ms.Extension

	article_file_b, err := os.ReadFile(this_article_path)

	if err != nil {

		return fmt.Errorf("failed to delete article: article file: %s", err.Error())
	}

	associatedTargets, err := GetAssociateMediaKeysForEditorjsSrc(article_file_b)

	if err != nil {

		return fmt.Errorf("failed to delete article: %s", err.Error())

	}

	atLen := len(associatedTargets)

	for i := 0; i < atLen; i++ {

		err := DeleteMedia(associatedTargets[i])

		if err != nil {

			return fmt.Errorf("failed to delete associated key: %s", err.Error())
		}

	}

	err = os.Remove(this_article_path)

	if err != nil {

		return fmt.Errorf("failed to delete article: rmart: %s", err.Error())
	}

	err = os.Remove(this_file_path)

	if err != nil {

		return fmt.Errorf("failed to delete video: rmmd: %s", err.Error())
	}

	return nil
}

func GetArticle(media_key string) (string, error) {

	var ms MediaStruct

	var content string

	this_file_path := mediaPath + media_key + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return "", fmt.Errorf("failed to get article: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &ms)

	if err != nil {

		return "", fmt.Errorf("failed to get article: marshal: %s", err.Error())
	}

	if ms.Type != "article" {

		return "", fmt.Errorf("failed to get article: %s: %s", "wrong type", ms.Type)

	}

	this_article_path := articlePath + media_key + "." + ms.Extension

	article_b, err := os.ReadFile(this_article_path)

	if err != nil {

		return "", fmt.Errorf("failed to get article: read file: %s", err.Error())

	}

	content = string(article_b)

	return content, nil

}

func UploadImage(c *gin.Context, author string, file *multipart.FileHeader, filename string, new_filename string, extension string) error {

	ms := MediaStruct{}

	this_file_path := mediaPath + new_filename + ".json"

	this_image_path := imagePath + new_filename + "." + extension

	c_time := time.Now()

	c_time_fmt := c_time.Format("2006-01-02-15-04-05")

	ms.ISPublic = true
	ms.Type = "image"
	ms.PlainName = filename
	ms.Extension = extension
	ms.CreatedAt = c_time_fmt
	ms.Author = author

	jb, err := json.Marshal(ms)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	err = os.WriteFile(this_file_path, jb, 0644)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	err = c.SaveUploadedFile(file, this_image_path)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	return nil
}

func DownloadMedia(c *gin.Context, watchId string) error {

	ms, err := GetByMediaKeyFromMedia(watchId)

	if ms == nil {

		return fmt.Errorf("failed to download media: %s", err.Error())

	}

	this_media_path := ""

	if ms.Type == "image" {

		this_media_path = imagePath + watchId + "." + ms.Extension

	} else if ms.Type == "video" {

		this_media_path = videoPath + watchId + "." + ms.Extension

	}

	if _, err := os.Stat(this_media_path); err != nil {

		return err

	}

	if ms.Type == "image" {

		c.Header("Content-Type", "image/"+ms.Extension)

	} else if ms.Type == "video" {

		c.Header("Content-Type", "video/"+ms.Extension)

	}

	c.File(this_media_path)

	return nil
}

func DownloadImage(c *gin.Context, watchId string) error {

	ms, err := GetByMediaKeyFromMedia(watchId)

	if ms == nil {

		return fmt.Errorf("failed to download image: %s", err.Error())

	}

	if ms.Type != "image" {

		return fmt.Errorf("failed to download image: %s: %s", "wrong type", ms.Type)
	}

	this_image_path := imagePath + watchId + "." + ms.Extension

	if _, err := os.Stat(this_image_path); err != nil {

		return err

	}

	c.Header("Content-Type", "image/"+ms.Extension)

	c.File(this_image_path)

	return nil
}

func UploadVideo(c *gin.Context, author string, file *multipart.FileHeader, filename string, new_filename string, extension string) error {

	ms := MediaStruct{}

	this_file_path := mediaPath + new_filename + ".json"

	this_video_path := videoPath + new_filename + "." + extension

	c_time := time.Now()

	c_time_fmt := c_time.Format("2006-01-02-15-04-05")

	ms.ISPublic = true
	ms.Type = "video"
	ms.PlainName = filename
	ms.Extension = extension
	ms.CreatedAt = c_time_fmt
	ms.Author = author

	jb, err := json.Marshal(ms)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	err = os.WriteFile(this_file_path, jb, 0644)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	err = c.SaveUploadedFile(file, this_video_path)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	return nil
}

func DeleteMedia(media_key string) error {

	var ms MediaStruct

	this_file_path := mediaPath + media_key + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return fmt.Errorf("failed to delete media: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &ms)

	if err != nil {

		return fmt.Errorf("failed to delete media: marshal: %s", err.Error())
	}

	var this_media_path string

	if ms.Type == "image" {

		this_media_path = imagePath + media_key + "." + ms.Extension

	} else if ms.Type == "video" {

		this_media_path = videoPath + media_key + "." + ms.Extension
	}

	err = os.Remove(this_media_path)

	if err != nil {

		return fmt.Errorf("failed to delete media: rm: %s", err.Error())
	}

	err = os.Remove(this_file_path)

	if err != nil {

		return fmt.Errorf("failed to delete media: rmkey: %s", err.Error())
	}

	return nil
}

func DeleteVideo(media_key string) error {

	var ms MediaStruct

	this_file_path := mediaPath + media_key + ".json"

	file_b, err := os.ReadFile(this_file_path)

	if err != nil {

		return fmt.Errorf("failed to delete video: %s", err.Error())

	}

	err = json.Unmarshal(file_b, &ms)

	if err != nil {

		return fmt.Errorf("failed to delete video: marshal: %s", err.Error())
	}

	this_video_path := videoPath + media_key + "." + ms.Extension

	err = os.Remove(this_video_path)

	if err != nil {

		return fmt.Errorf("failed to delete video: rmvid: %s", err.Error())
	}

	err = os.Remove(this_file_path)

	if err != nil {

		return fmt.Errorf("failed to delete video: rmmd: %s", err.Error())
	}

	return nil
}

func DownloadVideo(c *gin.Context, watchId string) error {

	ms, err := GetByMediaKeyFromMedia(watchId)

	if ms == nil {

		return fmt.Errorf("failed to download video: %s", err.Error())

	}

	if ms.Type != "video" {

		return fmt.Errorf("failed to download video: %s: %s", "wrong type", ms.Type)
	}

	this_video_path := videoPath + watchId + "." + ms.Extension

	if _, err := os.Stat(this_video_path); err != nil {

		return err

	}

	c.Header("Content-Type", "video/"+ms.Extension)

	c.File(this_video_path)

	return nil
}
