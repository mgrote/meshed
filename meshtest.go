package main

import (
	"log"
	"meshnode/dbclient"
	"meshnode/domain"
)

func main() {

	usernode := domain.NewUserNode("tralala", "hihi")

	user, ok := usernode.GetContent().(domain.User)
	if !ok {
		log.Fatal("could not convert user")
	}
	user.SetPassword("einszweidrei")
	usernode.SetContent(user)
	dbclient.Insert(usernode)

	secondnode := domain.NewUserNode("soso", "nanana")
	dbclient.Insert(secondnode)

	usernode.AddChild(secondnode)
	dbclient.Save(usernode)
	dbclient.Save(secondnode)



}
