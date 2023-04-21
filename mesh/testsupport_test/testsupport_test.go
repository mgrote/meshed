package testsupport_test

import (
	"github.com/franela/goblin"
	"github.com/mgrote/meshed/mesh/testsupport"
	"testing"
)

var caller = 0

func firstSuccessfullMethod() bool {
	caller++
	return true
}

func secondSuccessfullMethod() bool {
	caller++
	return true
}

func failingMethod() bool {
	caller++
	return false
}

func TestDoOnceFirstMethodCall(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Test First call to DoOnce", func() {
		g.It("DoOnce was not called", func() {
			g.Assert(caller).Equal(0)
		})
		g.It("Test First call to DoOnce", func() {
			g.Assert(testsupport.DoOnce("firstSuccessfullMethod", firstSuccessfullMethod)).Equal(true)
			g.Assert(caller).Equal(1)
		})
	})
	g.Describe("Test not ignoring other function call", func() {
		g.It("Not ignoring another method call to DoOnce", func() {
			// function should be ignored, if not fail here
			g.Assert(testsupport.DoOnce("secondSuccessfullMethod", secondSuccessfullMethod)).Equal(true)
			// function was not executed
			g.Assert(caller).Equal(2)
		})
	})
	g.Describe("Test failing function call", func() {
		g.It("Not ignoring failing function execution", func() {
			// function should be ignored, if not fail here
			g.Assert(testsupport.DoOnce("failingMethod", failingMethod)).Equal(false)
			// function was not executed
			g.Assert(caller).Equal(3)
		})
	})
	g.Describe("Test repeated failing function call", func() {
		g.It("Tries failing function execution again", func() {
			// function should be ignored, if not fail here
			g.Assert(testsupport.DoOnce("failingMethod", failingMethod)).Equal(false)
			// function was not executed
			g.Assert(caller).Equal(4)
		})
	})
}

func TestDoOnceSecondMethodCall(t *testing.T) {
	g := goblin.Goblin(t)
	// will fail if called separately, so do the first test first
	if caller == 0 {
		TestDoOnceFirstMethodCall(t)
	}
	g.Describe("Test ignoring function second call from test method with same key", func() {
		g.It("Ignoring following calls to DoOnce", func() {
			// function should be ignored, if not fail here
			g.Assert(testsupport.DoOnce("firstSuccessfullMethod", firstSuccessfullMethod)).Equal(true)
			// function was not executed
			g.Assert(caller).Equal(4)
		})
	})
}
