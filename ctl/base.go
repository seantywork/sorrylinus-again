package controller

import "github.com/gin-gonic/gin"

func GetIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{
		"title": "Index",
	})
}

func GetSigninIndex(c *gin.Context) {

	c.HTML(200, "signin.html", gin.H{
		"title": "Signin",
	})

}
