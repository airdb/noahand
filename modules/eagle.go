package modules

import (
	"github.com/astaxie/beego"
)

func Eagle() {
	for {
		beego.SetLogger("file", `{"filename":"/home/work/src/github.com/airdb/noah/logs/eagle.log"}`)
		beego.Info("eagle function")
	}
}
