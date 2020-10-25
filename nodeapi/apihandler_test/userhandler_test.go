package apihandler_test

import (
	"bytes"
	"encoding/json"
	"github.com/franela/goblin"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api user registeration", func() {
		registrationBody := map[string]interface{}{
			"user": "Heiner Müller",
			"pwd": "tralala-hihi",
			"email": "test@test.frup.de",
		}
		body, _ := json.Marshal(registrationBody)
		req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
		response := recordRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
	})
}

func TestLoginUser(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api user login", func() {
		loginBody := map[string]interface{}{
			"user": "Heiner Müller",
			"pwd": "tralala-hihi",
		}
		body, _ := json.Marshal(loginBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
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
