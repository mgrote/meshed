package apihandler_test

import (
	"context"
	"encoding/json"
	"github.com/franela/goblin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"meshed/configuration/configurations"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/model/categories"
	"meshed/meshnode/model/images"
	"meshed/meshnode/model/users"
	"meshed/meshnode/testsupport"
	"net/http"
	"testing"
	"time"
)

const meshDbTestConfigFile = "/Users/michaelgrote/etc/gotest/mesh.db.properties.ini"

func prepareTestDatabase() bool {
	dbclient.ReinitMeshDbClientWithConfig(meshDbTestConfigFile)
	dbConfig := configurations.ReadConfig(meshDbTestConfigFile)
	log.Println("testdatabase", dbConfig.Dbname, users.ClassName, "will be set empty")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collectionNames, error := dbclient.GridMongoClient.Database(dbConfig.Dbname).ListCollectionNames(ctx, bson.M{})
	if error != nil {
		log.Println("Got db error", error)
	}
	log.Println("Got collection names")
	for _, collectionName := range collectionNames {
		log.Println("Dropping collection", collectionName)
		dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(collectionName).Drop(ctx)
	}

	return true
}

func prepareTestData() bool {
	userNode := users.NewNode("Müller", "Heiner")
	user := users.GetUser(userNode)
	user.SetPassword("einszweidrei")
	userNode.SetContent(user)
	// save change content
	userNode.Save()

	secondUserNode := users.NewNode("Jelinek", "Elfriede")
	secondUser := users.GetUser(secondUserNode)
	secondUser.SetPassword("dreivier")
	secondUserNode.SetContent(secondUser)
	secondUserNode.Save()

	userImage := images.NewNode("user", "/Users/michaelgrote/Pictures/tusche/IMG_0294.jpeg")
	secondUserImage := images.NewNode("seconduser", "/Users/michaelgrote/Pictures/tusche/IMG_0311.jpeg")

	userNode.AddChild(userImage)
	secondUserNode.AddChild(secondUserImage)

	catOneNode := categories.NewNode("catone")
	catTwoNode := categories.NewNode("cattwo")

	catOneNode.AddChild(userImage)
	catTwoNode.AddChild(userImage)
	catTwoNode.AddChild(secondUserImage)

	userNode.AddChild(catOneNode)
	userNode.AddChild(catTwoNode)

	return true
}

func TestNodeTypes(t *testing.T) {
	testsupport.DoOnce("emptymeshdb", prepareTestDatabase)
	testsupport.DoOnce("storetestdata", prepareTestData)
	g := goblin.Goblin(t)
	g.Describe("Testing api listtypes", func() {
		req, _ := http.NewRequest("GET", "/listtypes", nil)
		response := recordRequest(req)
		g.It("Response code should be '200'/Http.OK", func() {
			g.Assert(response.Code).Equal(http.StatusOK)
		})
		g.It("Response should contain types", func() {
			var nodeTypes []string
			_ = json.Unmarshal(response.Body.Bytes(), &nodeTypes)
			g.Assert(containsString(nodeTypes, users.ClassName)).IsTrue()
			g.Assert(containsString(nodeTypes, images.ClassName)).IsTrue()
			g.Assert(containsString(nodeTypes, categories.ClassName)).IsTrue()
		})
	})
}

func containsString(a []string, e string) bool {
	for _, a := range a {
		if a == e {
			return true
		}
	}
	return false
}
