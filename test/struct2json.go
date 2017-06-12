package main

import (
	"encoding/json"
	"fmt"
	"github.com/airdb/noah/gvar"
)

func main() {
	// gvar.RunConfig.Comm.Interval = 600
	gvar.RunConfig.Comm.Listen = "127.0.0.1,localhost,innerip"
	gvar.RunConfig.Comm.Scheme = "http"
	gvar.RunConfig.Comm.Username = "apollo"
	gvar.RunConfig.Comm.Password = "3df077555fc40990d412237186c67a8"
	gvar.RunConfig.Comm.Uri = "/apollo/init.json"
	gvar.RunConfig.Comm.Server = "apollo.srehub.com"
	gvar.RunConfig.Comm.Server1 = "47.90.107.22"
	gvar.RunConfig.Comm.Server2 = "10.29.163.87"

	var module gvar.Module
	module.URL = "xxxxx`"

	for i := 1; i < 10; i++ {
		module.RunLevel = i
		gvar.RunConfig.ModuleList = append(gvar.RunConfig.ModuleList, module)
	}
	body, err := json.Marshal(gvar.RunConfig)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(body))
}
