package apihandler_test

import (
	"github.com/franela/goblin"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api user registeration", func() {
		req, _ := http.NewRequest("POST", "/register", nil)
		response := recordRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
	})
}

func TestLoginUser(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api user login", func() {
		req, _ := http.NewRequest("POST", "/login", nil)
		response := recordRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
	})
}

func TestRenewUserToken(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api renew user token", func() {
		req, _ := http.NewRequest("GET", "/renew", nil)
		response := recordRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
	})
}
