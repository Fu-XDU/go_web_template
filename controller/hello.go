package controller

import (
	"github.com/Fu-XDU/mingfu_go_common/base_response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, base_response.NewDataResponse("Hello World"))
}
