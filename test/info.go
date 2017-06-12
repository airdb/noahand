package main

import (
	"fmt"
	"github.com/airdb/noah/env"
	"github.com/airdb/sailor"
	"os"
	"runtime"
	"syscall"
	"time"
)

func main() {
	request := "hostname=" + env.Hostname() + "&" + "ostype=" + env.Ostype() + "&" + "Osrelease=" + env.Osrelease() + "&" + "timestamp=" + env.Timestamp()
	addr := "http://apollo.srehub.com/apollo/init.json?"
	requrl := addr + request

	for {
		ret, _ := sailor.HttpGet(requrl)
		fmt.Println(ret)
		time.Sleep(5 * time.Second)
	}
	select {}

	os.Exit(0)
	fmt.Println(requrl)

	fmt.Printf("OS: %s\nArchitecture: %s\n", runtime.GOOS, runtime.GOARCH)

	fmt.Println("GOPATH: ", runtime.GOROOT)
	if 0 == os.Getuid() {
		fmt.Println("getuid: ", os.Getuid())
	}

	PageSize := os.Getpagesize()
	fmt.Println("system pagesize: ", PageSize)

	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		fmt.Println(rlimit, err)
		return
	}
	fmt.Println("RLIMIT_NOFILE: ", rlimit.Cur)
	fmt.Println("RLIMIT_NOFILE: ", rlimit.Max)

	// 2. application ..
	GolangVersion := runtime.Version()
	fmt.Println("golang version: ", GolangVersion)

}
