package coremain

import (
	"fmt"
	"log"
	"plugin"
	"time"

	"github.com/airdb/noah/coremain/admin"
	"github.com/airdb/noah/internal/noahlib"
)

type Greeter interface {
	Greet()
}

func RunPlugin() {
	for {
		mod := noahlib.GetNoahPluginPath() + "plugin_greeter.so"

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

func Run() {
	admin.RunWeb()
	// Supervisor process.
	noahlib.StartSupervisor()

	// Heart beat, and get running config.
	go noahlib.RandomHeartbeat()

	// Worker process.
	for {
		log.Println("start_worker")
		// noahlib.RandomHeartbeat()
		noahlib.DownloadZip()
		// noahlib.DoSelfUpdate()
		log.Println("self update")
		time.Sleep(time.Minute)
	}
}
