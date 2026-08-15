package routes

import (
	"net/http"
	"path/filepath"
	"strings"

	mingfufile "github.com/Fu-XDU/mingfu_go_common/file"
	mingfuroutes "github.com/Fu-XDU/mingfu_go_common/routes"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

// 本地开发用 web/dist，Docker 镜像用 assets/web/v1
var webDistCandidates = []string{
	"./web/dist/",
	"./assets/web/v1/",
}

func resolveWebDistDir() string {
	for _, dir := range webDistCandidates {
		if ok, _ := mingfufile.FileExists(filepath.Join(dir, "index.html")); ok {
			return dir
		}
	}
	return webDistCandidates[0]
}

// 根路径重定向页面，使用 JavaScript 自动检测正确的跳转路径
const rootRedirectHTML = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Redirecting...</title></head>
<body>
<script>
// 获取当前路径，去掉末尾的斜杠
var base = window.location.pathname.replace(/\/$/, '');
// 跳转到 /v1/web/
window.location.replace(base + '/v1/web/');
</script>
<noscript><p>Redirecting... Please enable JavaScript.</p></noscript>
</body>
</html>`

func webRoutesV1(rg *gin.RouterGroup) {
	v1WebPath := resolveWebDistDir()
	log.Infof("serving web static files from %s at /v1/web/", v1WebPath)
	rg.Static("/web", v1WebPath)

	mingfuroutes.Router.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path

		// 根路径：返回 JS 重定向页面，自动检测 nginx 代理前缀
		if requestPath == "/" || requestPath == "" {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, rootRedirectHTML)
			return
		}

		// 检查是否是 /v1/web 路径下的请求（SPA 路由）
		// 只有这些请求才需要返回 index.html
		if strings.HasPrefix(requestPath, "/v1/web") {
			// 从请求路径中提取相对于 web 目录的路径
			// 例如: /v1/web/assets/index.js -> assets/index.js
			relativePath := strings.TrimPrefix(requestPath, "/v1/web")
			relativePath = strings.TrimPrefix(relativePath, "/")

			if relativePath != "" {
				file := filepath.Join(v1WebPath, relativePath)
				exist, _ := mingfufile.FileExists(file)
				if exist {
					c.File(file)
					return
				}
			}

			// 静态文件不存在，返回 index.html（SPA fallback）
			c.File(v1WebPath + "index.html")
			return
		}

		// 非 /v1/web 路径，返回 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}
