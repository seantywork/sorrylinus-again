package stream

import "github.com/gin-gonic/gin"

func CreateGenericServer() *gin.Engine {

	genserver := gin.Default()

	genserver.LoadHTMLGlob("view/*")

	genserver.Static("/public", "./public")

	return genserver

}
