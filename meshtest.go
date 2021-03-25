package main

import (
	"log"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/model/categories"
	"meshed/meshnode/model/images"
	"meshed/meshnode/model/users"
	"meshed/nodeapi/apirouting"
	"net/http"
)

func main() {

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

	firstUserImage := images.NewNode("user", "/Users/michaelgrote/Pictures/tusche/IMG_0294.jpeg")
	secondUserImage := images.NewNode("seconduser", "/Users/michaelgrote/Pictures/tusche/IMG_0311.jpeg")

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

	loaded, _ := dbclient.FindById(categories.ClassName, catOneNode.GetID())
	log.Println("got", loaded)

	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))
	// access api or close

}
