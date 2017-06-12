package env

import (
	"log"
	"os"
)

func Check() (flag bool) {
	if !IsRoot() {
		log.Println("the proc must be run as root, pls check!")
		os.Exit(0)
	}

	flag = true
	return

}
