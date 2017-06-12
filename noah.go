package main

import (
	"encoding/json"
	"fmt"
	"github.com/airdb/noah/env"
	"github.com/airdb/noah/gvar"
	"github.com/airdb/noah/modules"
	"github.com/airdb/sailor"
	"os"
	"strings"
	"time"
)

func Init() {
	os.Chdir(gvar.RunDir)
	env.Check()
	// Initialize RunConfig register noah agent, for updating runtime config from the remote apollo service.
	if sailor.IsEmptyValue(gvar.RunConfig.Comm.Interval) {
		gvar.RunConfig.Comm.Interval = 600
	}
	if sailor.IsEmptyValue(gvar.RunConfig.Comm.Listen) {
		gvar.RunConfig.Comm.Listen = "127.0.0.1,localhost,innerip"
	}
	if sailor.IsEmptyValue(gvar.RunConfig.Comm.Server) || sailor.IsEmptyValue(gvar.RunConfig.Comm.Scheme) || sailor.IsEmptyValue(gvar.RunConfig.Comm.Uri) {
		gvar.RunConfig.Comm.Interval = 5
		gvar.RunConfig.Comm.Scheme = "http"
		gvar.RunConfig.Comm.Username = "apollo"
		gvar.RunConfig.Comm.Password = "3df077555fc40990d412237186c67a8"
		gvar.RunConfig.Comm.Uri = "/apollo/init.json"
		gvar.RunConfig.Comm.Server = "apollo.srehub.com"
		gvar.RunConfig.Comm.Server1 = "47.90.107.22"
		gvar.RunConfig.Comm.Server2 = "10.29.163.87"
	}
	if gvar.RunConfigRollbackFlag {
		RunConfigRollback()
	}

	// fmt.Println(RunConfig.Comm.Server)
	// request param
	param := strings.Join([]string{env.Hostname(), env.Ostype(), env.Osrelease(), env.Timestamp()}, "&")

	// for domain in  RunConfig.Comm.Server
	// check domains avaliable or not.
	domain := gvar.RunConfig.Comm.Server

	requrl := strings.Join([]string{gvar.RunConfig.Comm.Scheme, "://", domain, gvar.RunConfig.Comm.Uri, "?", param}, "")
	fmt.Println(requrl)
	if !sailor.IsEmptyValue(gvar.RunConfig.Comm.Username) && !sailor.IsEmptyValue(gvar.RunConfig.Comm.Password) {
		requrl = strings.Join([]string{gvar.RunConfig.Comm.Scheme, "://", gvar.RunConfig.Comm.Username, ":", gvar.RunConfig.Comm.Password, "@", domain, gvar.RunConfig.Comm.Uri, "?", param}, "")
	}

	// http retry
	// gvar.NewConfigJson, _ = sailor.HttpGet(requrl)
	gvar.NewConfigJson, _ = sailor.HttpNewRequest("GET", requrl, nil)
	// fmt.Println(NewConfigJson)
	fmt.Println("apollo config size:", len(gvar.NewConfigJson))

	// New runtime Config
	ParseRunConfig()
	// Pseudo random interval time

	fmt.Println("apollo sleep: ", gvar.RunConfig.Comm.Interval)
	time.Sleep(time.Duration(gvar.RunConfig.Comm.Interval) * time.Second)
	// time.Sleep(10 * time.Second)
	// break
}

func ParseRunConfig() {
	err := json.Unmarshal([]byte(gvar.NewConfigJson), &gvar.NewRunConfig)
	if err == nil {
		gvar.LastRunConfig = gvar.RunConfig
		gvar.RunConfig = gvar.NewRunConfig

		// protect
		if gvar.RunConfig.Comm.Interval > gvar.RequstIntervalMax {
			gvar.RunConfig.Comm.Interval = gvar.RequstIntervalMax
		}
		if gvar.RunConfig.Comm.Interval < gvar.RequstIntervalMin {
			gvar.RunConfig.Comm.Interval = gvar.RequstIntervalMin
		}
	}
}

func RunConfigRollback() {
	gvar.RunConfig = gvar.LastRunConfig
	// report rollback
}

func main() {

	// run into daemon
	// modules.Daemon(0, 1)
	// go Init()

	for {
		Init()
		modules.ColumbiaRun()
	}

	select {}
}
