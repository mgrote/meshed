package apihandler_test

import (
	"github.com/franela/goblin"
	"github.com/mgrote/meshed/nodeapi/apirouting"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testrouter = apirouting.NewRouter()

func TestIndexRoot(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api index root", func() {
		req, _ := http.NewRequest("GET", "/", nil)
		response := recordRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
	})
}

func recordRequest(req *http.Request) *httptest.ResponseRecorder {
	log.Println("Requesting", req.URL)
	recordedResponse := httptest.NewRecorder()
	testrouter.ServeHTTP(recordedResponse, req)
	log.Println("Got response", recordedResponse.Body, "with http code", recordedResponse.Code)
	return recordedResponse
}
