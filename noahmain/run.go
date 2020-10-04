package noahmain

import (
	"fmt"
	"log"
	"math/rand"
	"plugin"
	"time"

	"github.com/airdb/noah/internal/noahlib"
)

type Greeter interface {
	Greet()
}

func Run() {
	go Heartbeat()

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

		time.Sleep(time.Duration(t) * time.Second)

		noahlib.Heartbeat()
	}
}
