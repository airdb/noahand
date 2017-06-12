package modules

import (
	"fmt"
	"github.com/airdb/noah/gvar"
	"github.com/airdb/sailor"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func Download() {
}

func Unpack() {
}

func DownloadAndUnpack(module gvar.Module) {
	// dlpath := gvar.WorkDir +  "cbd878db6ca670418634f145359d8dfe0"
	moduleWorkDir := filepath.Join("/", gvar.RunDir, gvar.RunModulesDir, module.Md5sum)

	dlpath := filepath.Join("/", gvar.RunDir, gvar.RunModulesDir, module.Md5sum, module.Name+"."+module.Filetype)
	fmt.Println("dlpath: ", dlpath)
	// pre

	if "bin" == module.Filetype {
		dlpath = filepath.Join("/", gvar.RunDir, gvar.RunModulesDir, module.Md5sum, module.Filetype, module.Name)
	}

	sailor.Download(module.URL, dlpath)

	if "bin" == module.Filetype {
		sailor.ChmodAddUserPerm(dlpath)
		destfilepath := filepath.Join("/", gvar.RunDir, gvar.RunModulesDir, module.Md5sum, module.Filetype, gvar.Control)
		sailor.MakeFileLink(dlpath, destfilepath)
	} else if "tgz" == module.Filetype || "tar.gz" == module.Filetype {
		// cmd := exec.Command("tar", "-C"+moduleWorkDir, "xvpf", dlpath)
		cmd := exec.Command("tar", "xvpf", dlpath, "-C"+moduleWorkDir)
		err := cmd.Run()
		if err == nil {
			fmt.Println("cmd run success.")
		} else {
			fmt.Println("cmd run err:", err)
		}
	}

	// post
	return
	os.Exit(0)
	StartProcess()
}

func StartProcess() {
	proc := "ping"
	pidfile := filepath.Join(gvar.RunDir, "run", proc+".pid")

	if sailor.IsExistFile(pidfile) {
		pidstring := sailor.ReadContentFromFile(pidfile)
		procfile := filepath.Join("/", "/proc", pidstring, "stat")
		// check pid running or not, return if running.
		fmt.Println(procfile)
		if sailor.IsExistFile(procfile) {
			return
		} else {
		}
	}

	filePath, _ := filepath.Abs("/bin/ping")
	args := append([]string{filePath}, "baidu.com")
	procAttr := new(os.ProcAttr)
	procAttr.Dir = "/noah"
	// procAttr.Sys.SysProcAttr = "/noah"
	// procAttr.Sys.Setsid = true
	procAttr.Env = []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin"}
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	if process, err := os.StartProcess("/bin/ping", args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s", proc, err.Error())
	} else {
		fmt.Printf("%s running as pid %d", proc, process.Pid)
		if sailor.WriteStringToFile(strconv.Itoa(process.Pid), pidfile) {
			fmt.Printf("%s write pid file success.", proc)
		}
	}

}
