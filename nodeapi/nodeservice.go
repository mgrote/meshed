package main

import (
	"flag"
	"fmt"
	"github.com/mgrote/meshed/mesh"
	"github.com/mgrote/meshed/nodeapi/apirouting"

	"log"
	"net/http"
	"os"
)

func main() {
	// read database location
	var pathFlag string
	flag.StringVar(&pathFlag, "inifiles", ".", "Path to ini files")
	flag.Parse()

	if err := mesh.InitApi(); err != nil {
		fmt.Println("could not initialize API")
		os.Exit(1)
	}

	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))
}
