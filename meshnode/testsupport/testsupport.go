package testsupport

import (
	"flag"
	"meshed/configuration/configurations"
)

type DoOnceFunction func() bool
var executedKeys = make([]string, 0)

func ReadFlags() {
	var pathFlag string
	flag.StringVar(&pathFlag, "inifiles", ".", "Path to ini files")
	flag.Parse()
	configurations.IniFilePath = pathFlag
}

func DoOnce(key string, doOnceUntilReset DoOnceFunction) bool {
	var success = true
	if !isAlreadyExecuted(key) {
		success = doOnceUntilReset()
		if success {
			executedKeys = append(executedKeys, key)
		}
	}
	return success
}

func isAlreadyExecuted(key string) bool {
	for _, executedKey :=  range executedKeys {
		if executedKey == key {
			return true
		}
	}
	return false
}

func Reset() {
	executedKeys = make([]string, 0)
}



