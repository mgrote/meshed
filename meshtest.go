package main

import (
	"meshnode/domain"
)

func main() {

	userNode := domain.NewUserNode("tralala", "hihi")

	user := domain.GetUser(userNode)
	user.SetPassword("einszweidrei")
	userNode.SetContent(user)
	userNode.Save()

	secondNode := domain.NewUserNode("soso", "nanana")
	secondUser := domain.GetUser(secondNode)
	secondUser.SetPassword("dreivier")
	secondNode.SetContent(secondUser)
	secondNode.Save()

	//userNode.AddChild(secondNode)
	secondNode.AddParent(userNode)

	userNode.RemoveChild(secondNode)


}
