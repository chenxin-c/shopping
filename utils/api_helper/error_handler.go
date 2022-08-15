package api_helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误处理
func HandleError(g *gin.Context, err error) {

	g.JSON(
		http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	g.Abort() //Abort 在被调用的函数中阻止挂起函数。注意这将不会停止当前的函数。例如，你有一个验证当前的请求是否是认证过的 Authorization 中间件。如果验证失败(例如，密码不匹配)，调用 Abort 以确保这个请求的其他函数不会被调用。
	return

}
