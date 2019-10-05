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

	userNode := users.NewNode("tralala", "hihi")

	user := users.GetUser(userNode)
	user.SetPassword("einszweidrei")
	userNode.SetContent(user)
	userNode.Save()

	secondNode := users.NewNode("soso", "nanana")
	secondUser := users.GetUser(secondNode)
	secondUser.SetPassword("dreivier")
	secondNode.SetContent(secondUser)
	secondNode.Save()

	//userNode.AddChild(secondNode)
//	secondNode.AddParent(userNode)
//	userNode.RemoveChild(secondNode)

	userImage := images.NewNode("user", "/Users/michaelgrote/Pictures/tusche/IMG_0294.jpeg")
	secondUserImage := images.NewNode("seconduser", "/Users/michaelgrote/Pictures/tusche/IMG_0311.jpeg")

	userNode.AddChild(userImage)
	secondNode.AddChild(secondUserImage)

	catOneNode := categories.NewNode("catone")
	catTwoNode := categories.NewNode("cattwo")

	catOneNode.AddChild(userImage)
	catTwoNode.AddChild(userImage)
	catTwoNode.AddChild(secondUserImage)

	for _, imageNode := range catTwoNode.GetChildren("image") {
		log.Println("got image node ", imageNode.GetContent())
		for _, userNode := range imageNode.GetParents("user") {
			log.Println("got user node", userNode.GetContent())
		}
	}

	loaded := dbclient.FindById(categories.ClassName, catOneNode.GetID())
	log.Println("got", loaded)

	router := apirouting.NewRouter()
	log.Fatal(http.ListenAndServe(":8001", router))


}
