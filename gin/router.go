package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/lhl1115/goencrypt/utils"
)

func NewRouter() *gin.Engine {

	r := gin.Default()
	r.Use(Encrypt(utils.AesKey, utils.AesIv, utils.AppID, utils.AppSecret, utils.ApiTimeout))

	r.GET("/ping", Ping)
	r.POST("/pong", Pong)

	return r
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ping",
	})
}
