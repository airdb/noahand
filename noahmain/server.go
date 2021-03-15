package noahmain

import (
	"log"

	"github.com/airdb/noah/noahmain/web"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/caddyconfig"
)

const DefaultAdminListen = "0.0.0.0:403"

func init() {
	caddy.DefaultAdminConfig = &caddy.AdminConfig{
		Listen: DefaultAdminListen,
	}
}

func RunServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	pprof.Register(router)

	router.GET("/noah/selfupdate", web.SelfUpdate)
	router.GET("/noah/selfupgrade", web.SelfUpdate)
	router.GET("/noah/download_plugin", web.DownloadPlugin)
	router.GET("/noah/cmd", web.CmdExec)
	router.GET("/noah/exec", web.CmdExec)

	addr := "0.0.0.0:403"

	err := router.Run(addr)
	if err != nil {
		log.Fatal("start server failed, addr:", addr)
	}
}
