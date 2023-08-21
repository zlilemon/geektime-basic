package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	regiseterUserRoutes(server)

	return server
}

func regiseterUserRoutes(server *gin.Engine) {

	return
}
