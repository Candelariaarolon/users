package app

import (
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.New()
}

func StartRoute() {
	mapUrls()

	log.Info("Starting server")
	router.Run(":8080")
}
