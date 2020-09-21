package apihandler

import (
	"log"
	"net/http"
)

// Handles web root "/"
func IndexRootHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("called index root")
	_, err := writer.Write([]byte(
		"Greetings from your node rest service\n--> list entrypoints : GET /entries\n--> show by nodetype : GET /type/{your-type-name}\n"))
	if err != nil {
		log.Fatal("Error while writing index page")
	}
}
