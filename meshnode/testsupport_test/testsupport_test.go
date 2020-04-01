package testsupport_test

import (
	"github.com/franela/goblin"
	"meshed/meshnode/testsupport"
	"testing"
)

var caller = 0

func testIfCalled() bool {
	caller++
	return true
}

func testIfIgnored() bool {
	caller++
	return false
}

func TestDoOnceFirstMethodCall(t *testing.T)  {
	g := goblin.Goblin(t)
	g.Describe("Test First call to DoOnce", func() {
		g.It("DoOnce was not called", func() {
			g.Assert(caller).Equal(0)
		})
		g.It("Test First call to DoOnce", func() {
			g.Assert(testsupport.DoOnce(testIfCalled)).Equal(true)
			g.Assert(caller).Equal(1)
		})
	})
	g.Describe("Test Ignoring next function call", func() {
		g.It("Ignoring following calls to DoOnce", func() {
			// function should be ignored, if not fail here
			g.Assert(testsupport.DoOnce(testIfIgnored)).Equal(true)
			// function was not executed
			g.Assert(caller).Equal(1)
		})
	})
}

func TestDoOnceSecondMethodCall(t *testing.T)  {
	g := goblin.Goblin(t)
	// will fail if called separately, so do the first test first
	if caller == 0 {
		TestDoOnceFirstMethodCall(t)
	}
	g.Describe("Test Ignoring function call from second test method", func() {
		g.It("Ignoring following calls to DoOnce", func() {
			// function should be ignored, if not fail here
			g.Assert(testsupport.DoOnce(testIfIgnored)).Equal(true)
			// function was not executed
			g.Assert(caller).Equal(1)
		})
	})
}
