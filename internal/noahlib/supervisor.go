package noahlib

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

func setProcessName(name string) {
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(unsafe.StringData(os.Args[0])))[:len(name)+1]

	n := copy(argv0, name+"\x00")
	if n < len(argv0) {
		argv0[n] = 0
	}
}

func StartSupervisor() {
	executable, _ := os.Executable()
	_ = os.Chdir(filepath.Base(executable))

	// deep copy os.Args & os.Environ
	osArgs := strings.Split(strings.Join(os.Args, "\x00"), "\x00")
	osEnviron := strings.Split(strings.Join(os.Environ(), "\x00")+"\x00master=0", "\x00")

	if os.Getenv("master") != "0" {
		setProcessName(filepath.Base(executable) + ":master process " + executable)
	} else {
		setProcessName(filepath.Base(executable) + ":worker process " + executable)

		return
	}

	var childs []*os.Process

	var supervisor func()

	supervisor = func() {
		p, err := os.StartProcess(executable, osArgs, &os.ProcAttr{
			Dir:   ".",
			Env:   osEnviron,
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if err != nil {
			log.Println(executable)

			return
		}

		for _, child := range childs {
			_ = child.Signal(syscall.SIGHUP)
		}

		childs = append(childs[:0], p)

		ps, err := p.Wait()
		if err != nil {
			log.Println(err.Error())
		}

		if ps != nil && !ps.Success() {
			go supervisor()
		}
	}

	go supervisor()
	go RandomHeartbeat()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT)
	signal.Notify(signalChan, syscall.SIGHUP)

	for {
		switch sig := <-signalChan; sig {
		case syscall.SIGTERM, syscall.SIGINT:
			for _, child := range childs {
				_ = child.Signal(sig)
			}

			os.Exit(0)
		case syscall.SIGHUP:
			log.Println("master start worker")

			go supervisor()
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

func WaitWorkerSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGHUP)

	switch <-c {
	case syscall.SIGTERM, syscall.SIGINT:
		log.Println("worker received exited signal")

		for _, f := range globalExitFuncs {
			f()
		}

		for _, f := range globalGraceFuncs {
			f()
		}

		return
	}

	executable, _ := os.Executable()
	setProcessName(filepath.Base(executable) + ": worker process (graceful shutdown)")
	log.Println("worker start graceful shutdown")

	for _, f := range globalGraceFuncs {
		f()
	}

	log.Println("worker end graceful shutdown")
}
