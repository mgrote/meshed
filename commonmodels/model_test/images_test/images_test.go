package images_test

import (
	"context"
	"github.com/franela/goblin"
	"github.com/mgrote/meshed/commonmodels/blobs"
	"github.com/mgrote/meshed/configurations"
	"github.com/mgrote/meshed/mesh/mongodb"
	"github.com/mgrote/meshed/mesh/testsupport"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
	mongodb.InitFileStreamDbClientWithConfig(gridDbConfigFile)
	dbConfig := configurations.ReadDbConfig(gridDbConfigFile)
	log.Println("testdatabase", dbConfig.DbName, dbConfig.ImageBucketName, "will be set empty")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongodb.GridMongoClient.Database(dbConfig.DbName).Collection(dbConfig.ImageBucketName + ".files").Drop(ctx)
	mongodb.GridMongoClient.Database(dbConfig.DbName).Collection(dbConfig.ImageBucketName + ".chunks").Drop(ctx)
	return true
}

func TestImageUpload(t *testing.T) {
	testsupport.DoOnce("emptymeshdb", prepareTestDatabase)
	g := goblin.Goblin(t)
	g.Describe("Blob upload should return an object id and no errors", func() {
		smallImageId, size, err := mongodb.UploadFile(smallImageFile, path.Base(smallImageFile))
		g.Assert(err == nil).IsTrue()
		g.Assert(size > 0).IsTrue()
		g.Assert(smallImageId != primitive.NilObjectID).IsTrue()
		largeImageId, size, err := mongodb.UploadFile(largeImageFile, path.Base(largeImageFile))
		g.Assert(err == nil).IsTrue()
		g.Assert(size > 0).IsTrue()
		g.Assert(largeImageId != primitive.NilObjectID).IsTrue()
		//dbclient.UploadFile(veryLargeImageFile , path.Base(veryLargeImageFile))
	})
}

func TestImageDownload(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Blob download", func() {
		log.Println("Download", path.Base(largeImageFile), "to", largeImageFileDownload)
		mongodb.DownloadFileByName(path.Base(largeImageFile), largeImageFileDownload)
		g.It("Downloaded file should exist in filesystem", func() {
			g.Assert(blobs.ReadableFile(largeImageFileDownload)).IsTrue()
		})
	})
}
