package images_test

import (
	"context"
	"github.com/franela/goblin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"meshed/configuration/configurations"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/model/images"
	"meshed/meshnode/testsupport"
	"os"
	"path"
	"testing"
	"time"
)

const gridDbConfigFile = "imagestream.db.properties.ini"

const smallImageFile = "/Users/michaelgrote/etc/gotest/particlestorm-16.jpg"
const largeImageFile = "/Users/michaelgrote/etc/gotest/PIA23623.jpg"
const veryLargeImageFile = "/Users/michaelgrote/etc/gotest/PIA23623_M34.tif"

const smallImageFileDownload = "/Users/michaelgrote/Downloads/particlestorm-16.jpg"
const largeImageFileDownload = "/Users/michaelgrote/Downloads/PIA23623.jpg"
const veryLargeImageFileDownload = "/Users/michaelgrote/Downloads/PIA23623_M34.tif"

func TestMain(m *testing.M) {
	testsupport.ReadFlags()
	os.Exit(m.Run())
}

func prepareTestDatabase() bool {
	dbclient.ReinitFileStreamDbClientWithConfig(gridDbConfigFile)
	dbConfig := configurations.ReadDbConfig(gridDbConfigFile)
	log.Println("testdatabase", dbConfig.Dbname, dbConfig.Bucketname, "will be set empty")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(dbConfig.Bucketname + ".files").Drop(ctx)
	dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(dbConfig.Bucketname + ".chunks").Drop(ctx)
	return true
}

func TestImageUpload(t *testing.T) {
	testsupport.DoOnce("emptymeshdb", prepareTestDatabase)
	g := goblin.Goblin(t)
	g.Describe("Image upload should return an object id and no errors", func() {
		smallImageId, size, err := dbclient.UploadFile(smallImageFile , path.Base(smallImageFile))
		g.Assert(err == nil).IsTrue()
		g.Assert(size > 0).IsTrue()
		g.Assert(smallImageId != primitive.NilObjectID).IsTrue()
		largeImageId, size, err :=dbclient.UploadFile(largeImageFile , path.Base(largeImageFile))
		g.Assert(err == nil).IsTrue()
		g.Assert(size > 0).IsTrue()
		g.Assert(largeImageId != primitive.NilObjectID).IsTrue()
		//dbclient.UploadFile(veryLargeImageFile , path.Base(veryLargeImageFile))
	})
}

func TestImageDownload(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Image download", func() {
		log.Println("Download", path.Base(largeImageFile), "to", largeImageFileDownload)
		dbclient.DownloadFileByName(path.Base(largeImageFile), largeImageFileDownload)
		g.It("Downloaded file should exist in filesystem", func() {
			g.Assert(images.ReadableFile(largeImageFileDownload )).IsTrue()
		})
	})
}
