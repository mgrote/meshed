package main

import (
	"log"
	"meshnode/model/categories"
	"meshnode/model/images"
	"meshnode/model/users"
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



}
