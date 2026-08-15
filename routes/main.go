package routes

import (
	mingfuflags "github.com/Fu-XDU/mingfu_go_common/flags"
	mingfuroutes "github.com/Fu-XDU/mingfu_go_common/routes"
	"github.com/gin-gonic/gin"
)

// Run will start the server
func Run() {
	mingfuroutes.Run(mingfuflags.ServerPort, mingfuflags.SslCertPath, mingfuflags.SslKeyPath, mingfuflags.TrustedProxies, getRoutes)
}

func getRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	webRoutesV1(v1)
	helloRoutesV1(v1)
}
