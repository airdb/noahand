/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/airdb/noah/noahmain"
	"github.com/airdb/noah/pkg/reap"
	"github.com/spf13/cobra"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func initRun() {
	rootCmd.AddCommand(runCmd)
}

func Run() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Set to 0 because we're doing our own time, with timezone

	if os.Getenv("_WORKER") == "true" {
		noahmain.Run()
	} else {
		RunMaster()
	}
}

func RunMaster() {
	// Reap all children worker with no control or feedback.
	go reap.ReapChildren(nil, nil, nil, nil)

	var pids []int

	// Ensure all subprocesses are killed.
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		for range sigchan {
			for _, pid := range pids {
				log.Printf("kill process, pid=%v", pid)

				err := syscall.Kill(pid, syscall.SIGKILL)
				if err != nil {
					continue
				}

				time.Sleep(time.Second)
			}
		}

		os.Exit(0)
	}()

	pids = StartWorker()

	if len(pids) == 0 {
		log.Fatalf("launch workers failed, args=%v", os.Args)
	}

	for {
		var alivePids []int

		for _, pid := range pids {
			process, err := os.FindProcess(pid)
			if err != nil {
				log.Printf("find process failed, pid=%v", pid)

				continue
			}

			err = process.Signal(syscall.Signal(0))
			if err != nil {
				log.Printf("find process failed, pid=%v,err=%v", pid, err)

				continue
			}

			alivePids = append(alivePids, pid)
		}

		if len(alivePids) > 0 {
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

const (
	RoleWorker = "eagle"
	RoleMaster = "master"
)

func StartWorker() (pids []int) {
	sWorkerNum, present := os.LookupEnv("WORKER_NUM")
	if !present {
		sWorkerNum = "1"
		log.Printf("WORKER_NUM is null, set default=%v", sWorkerNum)
	}

	workerNum, err := strconv.Atoi(sWorkerNum)
	if err != nil {
		log.Fatalf("parse_worker_num_failed, worker_num=%v,err=%v", sWorkerNum, err)
	}

	for i := 0; i < workerNum; i++ {
		pid, err := RunWorker(i)
		if err != nil {
			log.Printf("start_worker_failed, args=%v,err=%v", os.Args, err)

			continue
		}

		pids = append(pids, pid)
	}

	return pids
}

func RunWorker(id int) (pid int, err error) {
	// Set worker specific environment variables.
	envs := os.Environ()
	envs = append(envs, "_WORKER=true")
	envs = append(envs, fmt.Sprintf("_WORKER_ID=%d", id))

	execSpec := &syscall.ProcAttr{
		Env: envs,
		Sys: &syscall.SysProcAttr{Setpgid: true},
		Files: []uintptr{
			os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(),
		},
	}

	if os.Args[len(os.Args)-1] == RoleMaster {
		os.Args[len(os.Args)-1] = RoleWorker
	}

	if os.Args[len(os.Args)-1] != RoleWorker {
		os.Args = append(os.Args, RoleWorker)
	}

	// nolint: gosec
	return syscall.ForkExec(os.Args[0], os.Args, execSpec)
}
