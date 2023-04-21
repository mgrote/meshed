package testsupport

import (
	"flag"
	"github.com/mgrote/meshed/configurations"
	"log"
)

type DoOnceFunction func() bool

var executedKeys = make([]string, 0)

func ReadFlags() {
	log.Println("Testsupport: checking flags")
	var pathFlag string
	flag.StringVar(&pathFlag, "inifiles", ".", "Path to ini files")
	flag.Parse()
	log.Println("Path to ini files", pathFlag)
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
	for _, executedKey := range executedKeys {
		if executedKey == key {
			return true
		}
	}
	return false
}

func Reset() {
	executedKeys = make([]string, 0)
}
