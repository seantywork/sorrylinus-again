package auth

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var DEBUG bool = false

func OauthGoogleLogin(c *gin.Context) {

	oauth_state := GenerateStateAuthCookie(c)

	u := GoogleOauthConfig.AuthCodeURL(oauth_state)

	c.Redirect(302, u)

}

func OauthGoogleCallback(c *gin.Context) {

	session := sessions.Default(c)

	var session_id string

	v := session.Get("SOLIAGAIN")

	if v == nil {
		fmt.Printf("access auth failed: %s\n", "session id not found")
		return
	} else {
		session_id = v.(string)
	}

	state := c.Request.FormValue("state")

	if state == "" {
		fmt.Printf("access auth failed: %s\n", "form state not found")
		return
	}

	if state != session_id {
		fmt.Printf("access auth failed: %s\n", "value not matching")
		c.Redirect(302, "/signin")
		return
	}

	data, err := GetUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	var oauth_struct OAuthStruct

	err = json.Unmarshal(data, &oauth_struct)

	if err != nil {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	if !oauth_struct.VERIFIED_EMAIL {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	if err != nil {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	fmt.Println("oauth sign in success")

	c.Redirect(302, "/")

}
