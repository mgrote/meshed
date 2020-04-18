package users_test

import (
	"context"
	"fmt"
	"github.com/franela/goblin"
	"meshed/configuration/configurations"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/model/users"
	"meshed/meshnode/testsupport"
	"reflect"
	"testing"
	"time"
)

const meshDbConfigFile = "/Users/michaelgrote/etc/gotest/mesh.db.properties.ini"

func prepareTestDatabase() bool {
	dbclient.ReinitMeshDbClientWithConfig(meshDbConfigFile)
	dbConfig := configurations.ReadConfig(meshDbConfigFile)
	fmt.Println("testdatabase", dbConfig.Dbname, users.ClassName, "will be set empty")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(users.ClassName).Drop(ctx)
	return true
}

func TestUserCreation(t *testing.T)  {
	testsupport.DoOnce(prepareTestDatabase)
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
	testsupport.DoOnce(prepareTestDatabase)
	g := goblin.Goblin(t)
	g.Describe("User password", func() {
		userNode := users.NewNode("Hüter", "Horst der")
		user := users.GetUser(userNode)
		user.SetPassword("onetwothree")
		userNode.SaveContent(user)
		g.It("Password should be encrypted", func() {
			g.Assert(user.Password).IsNotEqual("onetwothree")
		})
		g.It("Password should be approved", func() {
			g.Assert(user.IsPassword("onetwothree")).IsTrue()
		})
	})
}
