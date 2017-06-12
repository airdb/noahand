package gvar

import (
	"unsafe"
)

const (
	// VERSION represent beego web framework version.
	VERSION = "1.0.0"
	// DEV is for develop
	DEV = "dev"
	// PROD is for production
	PROD = "prod"
)

// global config value, Real-time update
var LastRunConfig RuntimeConfig
var RunConfig RuntimeConfig
var NewRunConfig RuntimeConfig
var LastConfigJson string
var RunConfigJson string
var NewConfigJson string

var DefaultRunConfigSize = unsafe.Sizeof(RunConfig)
var RunConfigRollbackFlag bool

//  /noah/modules/* if run user id is 0,  $HOME/noah/modules/* if run user id not 0.
var RunDir = "noah"
var RunModulesDir = "modules"
var Control = "control"

var RequstIntervalMax = 600
var RequstIntervalMin = 5
