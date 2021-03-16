package coremain

import (
	"log"

	"airdb.io/airdb/noah/coremain/web"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/caddyserver/caddy/v2"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"

)

const DefaultAdminListen = "0.0.0.0:403"

func init() {
	caddy.DefaultAdminConfig = &caddy.AdminConfig{
		Listen: DefaultAdminListen,
	}
}

func RunA() {
	caddycmd.Main()
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
