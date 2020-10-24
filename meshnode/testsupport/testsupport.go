package testsupport

type DoOnceFunction func() bool
var executedKeys = make([]string, 0)

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



