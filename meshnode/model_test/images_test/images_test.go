package images_test

import (
	"context"
	"fmt"
	"github.com/franela/goblin"
	"meshed/configuration/configurations"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/model/images"
	"meshed/meshnode/testsupport"
	"path"
	"testing"
	"time"
)

const gridDbConfigFile = "/Users/michaelgrote/etc/gotest/imagestream.db.properties.ini"

const smallImageFile = "/Users/michaelgrote/etc/gotest/particlestorm-16.jpg"
const largeImageFile = "/Users/michaelgrote/etc/gotest/PIA23623.jpg"
const veryLargeImageFile = "/Users/michaelgrote/etc/gotest/PIA23623_M34.tif"

const smallImageFileDownload = "/Users/michaelgrote/Downloads/particlestorm-16.jpg"
const largeImageFileDownload = "/Users/michaelgrote/Downloads/PIA23623.jpg"
const veryLargeImageFileDownload = "/Users/michaelgrote/Downloads/PIA23623_M34.tif"

func prepareTestDatabase() bool {
	dbclient.ReinitFileStreamDbClientWithConfig(gridDbConfigFile)
	dbConfig := configurations.ReadConfig(gridDbConfigFile)
	fmt.Println("testdatabase", dbConfig.Dbname, dbConfig.Bucketname, "will be set empty")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(dbConfig.Bucketname + ".files").Drop(ctx)
	dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(dbConfig.Bucketname + ".chunks").Drop(ctx)
	return true
}

func TestImageUpload(t *testing.T) {
	testsupport.DoOnce(prepareTestDatabase)
	g := goblin.Goblin(t)
	g.Describe("Image upload", func() {
		dbclient.UploadFile(smallImageFile , path.Base(smallImageFile))
		dbclient.UploadFile(largeImageFile , path.Base(largeImageFile))
		//dbclient.UploadFile(veryLargeImageFile , path.Base(veryLargeImageFile))
	})
}

func TestImageDownload(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Image download", func() {
		fmt.Println("Download", path.Base(smallImageFile), "to", smallImageFileDownload)
		dbclient.DownloadFile(path.Base(smallImageFile), smallImageFileDownload)
		g.It("Downloaded file should exist in filesystem", func() {
			g.Assert(images.ReadableFile(smallImageFileDownload )).Equal(true)
		})
	})
}
