package modules

import (
	"fmt"
	"github.com/airdb/noah/gvar"
)

func ColumbiaRun() {
	if !gvar.RunConfigRollbackFlag {
		fmt.Println("Columbia is running, start eagle...")

		// xxxxx
		for _, module := range gvar.RunConfig.ModuleList {
			fmt.Println("Columbia is running...", module.RunLevel)
			DownloadAndUnpack(module)
			// if 1 == module.RunLevel {
			// }
		}

	}
}
