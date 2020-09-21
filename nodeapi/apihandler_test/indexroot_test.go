package apihandler_test

import (
	"github.com/franela/goblin"
	"log"
	"meshed/nodeapi/apirouting"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testrouter = apirouting.NewRouter()


func TestIndexRoot(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Testing api index root", func() {
		req, _ := http.NewRequest("GET", "/", nil)
		response := executeRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
		checkResponseCode(t, http.StatusOK, response.Code)
	})
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	log.Println("Requesting", req.URL)
	recordedResponse := httptest.NewRecorder()
	testrouter.ServeHTTP(recordedResponse, req)
	log.Println("Got response", recordedResponse.Body, "with http code", recordedResponse.Code)
	return recordedResponse
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
