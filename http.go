package fast

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Ghttp *gin.Engine

func RunHttp(rs ...Router) {
	Ghttp = gin.New()

	level := Conf.DefInt("log.level", 0)

	mode := gin.DebugMode
	if level != 0 {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	Ghttp.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "CONNECT", "TRACE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Token", "Authorization"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
			AllowAllOrigins:  true,
		}))

	parseRouter(Ghttp, rs...)

	err := Ghttp.Run(fmt.Sprintf(":%d", Conf.DefInt("port", 8080)))
	if err != nil {
		L.Fatalf("http服务启动失败(ERROR:%s)", err)
	}
}
