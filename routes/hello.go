package routes

import (
	"go_web_template/controller"
	"github.com/gin-gonic/gin"
)

func helloRoutesV1(rg *gin.RouterGroup) {
	{
		group := rg.Group("/hello")
		group.GET("/", controller.HelloWorld)
	}
}