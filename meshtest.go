package main

import (
	"flag"
	"fmt"
	"github.com/mgrote/meshed/commonmodels/blobs"
	"github.com/mgrote/meshed/commonmodels/categories"
	"github.com/mgrote/meshed/commonmodels/users"
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

	if err := mesh.InitApiWithConfig(pathFlag); err != nil {
		fmt.Println("init mesh api:", err)
		os.Exit(1)
	}

	// create some node with content
	firstUserNode := users.NewNode("User", "One")
	firstUser := users.GetUser(firstUserNode)
	firstUser.SetPassword("einszweidrei")
	firstUserNode.SetContent(firstUser)
	firstUserNode.Save()

	secondUserNode := users.NewNode("Other", "User")
	secondUser := users.GetUser(secondUserNode)
	secondUser.SetPassword("dreivier")
	secondUserNode.SetContent(secondUser)
	secondUserNode.Save()

	firstUserImage := blobs.NewNode("user", "/Users/michaelgrote/Pictures/tusche/IMG_0294.jpeg")
	secondUserImage := blobs.NewNode("seconduser", "/Users/michaelgrote/Pictures/tusche/IMG_0311.jpeg")

	firstUserNode.AddChild(firstUserImage)
	secondUserNode.AddChild(secondUserImage)

	catOneNode := categories.NewNode("catone")
	catTwoNode := categories.NewNode("cattwo")

	catOneNode.AddChild(firstUserImage)
	catTwoNode.AddChild(firstUserImage)
	catTwoNode.AddChild(secondUserImage)

	for _, imageNode := range catTwoNode.GetChildren("image") {
		log.Println("got image node ", imageNode.GetContent())
		for _, userNode := range imageNode.GetParents("user") {
			log.Println("got user node", userNode.GetContent())
		}
	}

	loaded, _ := mesh.Service.FindNodeById(categories.TypeName, catOneNode.GetID())
	log.Println("got", loaded)

	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))
	// access api or close

}
