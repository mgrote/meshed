package users_test

import (
	"github.com/franela/goblin"
	"meshed/meshnode/model/users"
	"reflect"
	"testing"
)

func TestUserCreation(t *testing.T)  {
	g := goblin.Goblin(t)
	g.Describe("User creation", func() {
		userNode := users.NewNode("Müller", "Heiner")
		//reflect.TypeOf(userContent).String()
		g.It("Node should have user as content", func() {
			g.Assert(reflect.TypeOf(userNode.GetContent()).String()).Equal("users.User")
		})
		user := users.GetUser(userNode)
		g.It("Should has name", func() {
			g.Assert(user.Forename).Equal("Heiner")
		})
	})
}

func TestUserPassword(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("User password", func() {
		userNode := users.NewNode("Hüter", "Horst")
		user := users.GetUser(userNode)
		user.SetPassword("onetwothree")
		g.It("Password should be encrypted", func() {
			g.Assert(user.Password).IsNotEqual("onetwothree")
		})
		g.It("Password should be approved", func() {
			g.Assert(user.IsPassword("onetwothree")).IsTrue()
		})
	})
}