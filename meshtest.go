package main

import (
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
	secondNode.AddParent(userNode)

	userNode.RemoveChild(secondNode)


}
