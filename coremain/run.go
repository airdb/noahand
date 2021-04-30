package coremain

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"plugin"
	"strconv"
	"strings"
	"syscall"
	"time"

	"airdb.io/airdb/noah/internal/noahlib"
	"airdb.io/airdb/noah/pkg/reap"
)

type Greeter interface {
	Greet()
}

func RunPlugin() {
	go Heartbeat()

	return
	for {
		mod := noahlib.GetPluginPath() + "plugin_greeter.so"

		plug, err := plugin.Open(mod)
		if err != nil {
			log.Println("open plugin failed", err)

			continue
		}

		// 2. look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		symGreeter, err := plug.Lookup("Greeter")
		if err != nil {
			fmt.Println(err)

			continue
		}

		// 3. Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		var greeter Greeter

		greeter, ok := symGreeter.(Greeter)
		if !ok {
			fmt.Println("unexpected type from module symbol")

			continue
		}

		// 4. use the module
		greeter.Greet()

		fmt.Println("run plugin...")
		time.Sleep(time.Minute)
	}
}

func Heartbeat() {
	for {
		// rand.Seed(900)
		maxSleepInterval := 60

		// nolint: gosec
		t := rand.Intn(maxSleepInterval)

		log.Println("sleep time seed", t)
		t =10

		time.Sleep(time.Duration(t) * time.Second)

		noahlib.Heartbeat()
	}
}

func Run() {
	// log.SetOutput(os.Stdout)
	// log.SetFlags(0) // Set to 0 because we're doing our own time, with timezone
	noahlib.StartSupervisor()

	Heartbeat()
}

func waitForSignals(pids []int) {
	sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGQUIT)

	for range sigchan {
		s := <-sigchan
		log.Printf("%v signal received.\n", s)

		for _, pid := range pids {
			switch s {
			case syscall.SIGHUP:
				log.Println("SIGHUP")
			case syscall.SIGUSR2:
				log.Println("SIGHUP")

				err := syscall.Kill(pid, syscall.SIGKILL)
				if err != nil {
					log.Println(err)

					continue
				}

				time.Sleep(time.Second)

				// if os.Getenv("_WORKER") != "true" {
				// 	RunMaster()
				// }
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				log.Println("quit")
				// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				// defer cancel()
				// return server.Shutdown(ctx)
				log.Printf("kill process, pid=%v", pid)

				err := syscall.Kill(pid, syscall.SIGKILL)
				if err != nil {
					continue
				}

				time.Sleep(time.Second)
			default:
				log.Printf("kill process, pid=%v", pid)

				err := syscall.Kill(pid, syscall.SIGKILL)
				if err != nil {
					continue
				}

				time.Sleep(time.Second)
			}
		}
	}
}

func RunMaster() {
	// Reap all children worker with no control or feedback.
	go reap.ReapChildren(nil, nil, nil, nil)

	var pids []int

	// Ensure all subprocesses are killed.
	go func() {
		waitForSignals(pids)
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

func Runa() {
	execSpec := &syscall.ProcAttr{
		Env: os.Environ(),
		Sys: &syscall.SysProcAttr{Setpgid: true},
		Files: []uintptr{
			os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(),
		},
	}

	for _, role := range []string{RoleMaster, RoleWorker} {
		args := append(os.Args, role)

		cmd := strings.Join(args, " ")

		isRunning, err := noahlib.CheckProRunning(cmd)
		fmt.Println(isRunning, err)
		// if err != nil {
		// 	return
		// }

		if isRunning {
			log.Printf("%v is running\n", os.Args)
			return
		}

		log.Printf("run '%v'\n", os.Args)
		syscall.ForkExec(os.Args[0], os.Args, execSpec)
	}
}
