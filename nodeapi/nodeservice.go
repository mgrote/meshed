package main

import (
	"flag"
	"log"
	"meshed/configuration/configurations"
	"net/http"

	"meshed/nodeapi/apirouting"
)

func main() {
	// read database location
	var pathFlag string
	flag.StringVar(&pathFlag, "inifiles", ".", "Path to ini files")
	flag.Parse()
	configurations.IniFilePath = pathFlag

	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))
}
