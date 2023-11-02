package main

import (
	"flag"
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/commonmodels/blobs"
	"github.com/mgrote/meshed/commonmodels/categories"
	"github.com/mgrote/meshed/commonmodels/users"
	"github.com/mgrote/meshed/mesh"
	"github.com/mgrote/meshed/meshserviceprovider/inmemorymap"
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

	// Init API with default config.
	//if err := mongodb.InitApiWithConfig(pathFlag); err != nil {
	//	fmt.Println("init mesh api:", err)
	//	os.Exit(1)
	//}
	if err := inmemorymap.InitApi(); err != nil {
		log.Println("init mesh api:", err)
		os.Exit(1)
	}

	// create some node with content
	firstUserNode := users.NewNode("User", "One")
	firstUser := users.GetUser(firstUserNode)
	firstUser.SetPassword("einszweidrei")
	firstUserNode.SetContent(firstUser)
	if err := firstUserNode.Save(); err != nil {
		os.Exit(3)
	}

	secondUserNode := users.NewNode("Other", "User")
	secondUser := users.GetUser(secondUserNode)
	secondUser.SetPassword("dreivier")
	secondUserNode.SetContent(secondUser)
	if err := secondUserNode.Save(); err != nil {
		os.Exit(4)
	}

	firstUserImage := blobs.NewNode("user", "/Users/michaelgrote/Pictures/tusche/IMG_0294.jpeg")
	secondUserImage := blobs.NewNode("seconduser", "/Users/michaelgrote/Pictures/tusche/IMG_0311.jpeg")

	if err := firstUserNode.AddChild(firstUserImage); err != nil {
		os.Exit(5)
	}
	if err := secondUserNode.AddChild(secondUserImage); err != nil {
		os.Exit(6)
	}

	catOneNode := categories.NewNode("catone")
	catTwoNode := categories.NewNode("cattwo")

	if err := catOneNode.AddChild(firstUserImage); err != nil {
		os.Exit(7)
	}
	if err := catTwoNode.AddChild(firstUserImage); err != nil {
		os.Exit(8)
	}
	if err := catTwoNode.AddChild(secondUserImage); err != nil {
		os.Exit(9)
	}

	for _, imageNode := range catTwoNode.GetChildren("image") {
		log.Println("got image node ", imageNode.GetContent())
		for _, userNode := range imageNode.GetParents("user") {
			log.Println("got user node", userNode.GetContent())
		}
	}

	loaded, _ := mesh.NodeService.FindNodeByID(commonmodels.CategoryType, catOneNode.GetID())
	log.Println("got", loaded)

	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))
	// access api or close

}
