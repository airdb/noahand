package web

import (
	"net/http"

	"airdb.io/airdb/noah/internal/noahlib"
	"airdb.io/airdb/sailor"
	"github.com/gin-gonic/gin"
)

func SelfUpdate(c *gin.Context) {
	noahlib.DoSelfUpdate()
	c.String(http.StatusOK, "upgrade successfully")
}

func DownloadPlugin(c *gin.Context) {
	noahlib.Downloader()
	c.String(http.StatusOK, "upgrade plugin successfully")
}

type CmdReq struct {
	Cmd  string   `form:"cmd"`
	Args []string `form:"args"`
}

func CmdExec(c *gin.Context) {
	var req CmdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusForbidden, "cmd not exists")

		return
	}

	ret, err := sailor.ExecCommand(req.Cmd, req.Args)
	if err != nil {
		c.String(http.StatusOK, "exec command failed")

		return
	}

	c.String(http.StatusOK, ret)
}
