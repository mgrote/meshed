package main

import (
	"log"
	"net/http"

	"meshed/nodeapi/apirouting"
)

func main() {
	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))
}
