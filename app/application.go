package app

import (
	"github.com/cesarnono/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplicacion() {
	mapUrls()
	logger.Info("Application initialized ......")
	router.Run(":8082")

}
