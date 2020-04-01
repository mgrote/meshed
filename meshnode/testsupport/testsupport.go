package testsupport

type DoOnceFunction func() bool
var alreadyExcecuted bool

func DoOnce(doOnceUntilReset DoOnceFunction) bool {
	if !alreadyExcecuted {
		alreadyExcecuted = doOnceUntilReset()
	}
	return alreadyExcecuted
}

func Reset() {
	alreadyExcecuted = false
}



