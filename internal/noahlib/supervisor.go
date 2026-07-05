package noahlib

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// setProcessName rewrites the memory backing argv[0] so the new name shows up
// in ps/top. The kernel-provided buffer is only len(os.Args[0])+1 bytes, so a
// longer name is truncated; writing past it would corrupt argv[1:]/environ.
func setProcessName(name string) {
	argv0 := (*[1 << 20]byte)(unsafe.Pointer(unsafe.StringData(os.Args[0])))[: len(os.Args[0])+1 : len(os.Args[0])+1]

	n := copy(argv0[:len(argv0)-1], name)
	for i := n; i < len(argv0); i++ {
		argv0[i] = 0
	}
}

func StartSupervisor() {
	executable, _ := os.Executable()
	_ = os.Chdir(filepath.Dir(executable))

	// Deep copy os.Args and os.Environ: their strings point into the original
	// argv/environ block, which setProcessName mutates below.
	osArgs := strings.Split(strings.Join(os.Args, "\x00"), "\x00")
	osEnviron := strings.Split(strings.Join(os.Environ(), "\x00")+"\x00master=0", "\x00")

	if os.Getenv("master") == "0" {
		setProcessName(filepath.Base(executable) + ":worker process - zeus")
		log.Println("Worker process started")
		return
	}

	setProcessName(filepath.Base(executable) + ":master process - eagle")

	type workerExit struct {
		proc    *os.Process
		success bool
	}
	exited := make(chan workerExit, 8)

	// child is only accessed from the supervisor loop below, so no locking.
	var child *os.Process

	startWorker := func() {
		log.Println("Starting worker process")
		p, err := os.StartProcess(executable, osArgs, &os.ProcAttr{
			Dir:   ".",
			Env:   osEnviron,
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if err != nil {
			log.Println("Error starting process:", err)
			return
		}

		old := child
		child = p
		if old != nil {
			_ = old.Signal(syscall.SIGHUP)
		}

		go func() {
			ps, err := p.Wait()
			if err != nil {
				log.Println("Process wait error:", err)
			}
			exited <- workerExit{proc: p, success: ps != nil && ps.Success()}
		}()
	}

	startWorker()
	go RandomHeartbeat()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	for {
		select {
		case e := <-exited:
			if e.proc != child {
				// An old worker replaced during a SIGHUP reload; nothing to do.
				continue
			}
			child = nil
			if e.success {
				log.Println("Worker process exited normally, not restarting")
				continue
			}
			log.Println("Worker process exited with failure, restarting")
			time.Sleep(time.Second) // avoid a tight fork loop when the worker crashes at startup
			startWorker()
		case sig := <-signalChan:
			switch sig {
			case syscall.SIGHUP:
				log.Println("Received SIGHUP, starting worker")
				startWorker()
			default:
				log.Println("Received signal:", sig)
				if child != nil {
					_ = child.Signal(sig)
					timeout := time.After(10 * time.Second)
					for waiting := true; waiting; {
						select {
						case e := <-exited:
							waiting = e.proc != child
						case <-timeout:
							log.Println("Worker did not exit in time, killing")
							_ = child.Kill()
							waiting = false
						}
					}
				}
				os.Exit(0)
			}
		}
	}
}

var (
	globalMu         sync.Mutex
	globalExitFuncs  []func()
	globalGraceFuncs []func()
)

func RegisterExitSignals(f func()) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalExitFuncs = append(globalExitFuncs, f)
}

func RegisterGraceSignals(f func()) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalGraceFuncs = append(globalGraceFuncs, f)
}

func registeredFuncs() (exitFuncs, graceFuncs []func()) {
	globalMu.Lock()
	defer globalMu.Unlock()

	return append([]func(){}, globalExitFuncs...), append([]func(){}, globalGraceFuncs...)
}

func WaitWorkerSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	sig := <-c
	exitFuncs, graceFuncs := registeredFuncs()

	switch sig {
	case syscall.SIGTERM, syscall.SIGINT:
		log.Println("worker received exit signal")

		for _, f := range exitFuncs {
			f()
		}

		for _, f := range graceFuncs {
			f()
		}
	case syscall.SIGHUP:
		executable, _ := os.Executable()
		setProcessName(filepath.Base(executable) + ":worker (grace)")
		log.Println("worker start graceful shutdown")

		for _, f := range graceFuncs {
			f()
		}

		log.Println("worker end graceful shutdown")
	}
}
