package apihandler_test

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"meshed/configuration/configurations"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/testsupport"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const gridDbTestConfigFile = "imagestream.db.properties.ini"

func TestMain(m *testing.M) {
	testsupport.ReadFlags()
	os.Exit(m.Run())
}

func prepareImageTestDatabase() bool {
	dbclient.ReinitFileStreamDbClientWithConfig(gridDbTestConfigFile)
	dbConfig := configurations.ReadDbConfig(gridDbTestConfigFile)
	fmt.Println("testdatabase", dbConfig.Dbname, dbConfig.Bucketname, "will be set empty")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(dbConfig.Bucketname + ".files").Drop(ctx)
	err = dbclient.GridMongoClient.Database(dbConfig.Dbname).Collection(dbConfig.Bucketname + ".chunks").Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// taken from
// https://stackoverflow.com/questions/43904974/testing-go-http-request-formfile
func TestUploadImage(t *testing.T) {
	testsupport.DoOnce("emptyimagedb", prepareImageTestDatabase)
	testsupport.DoOnce("emptymeshdb", prepareTestDatabase)
	//Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	//This writers is going to transform
	//what we pass to it to multipart form data
	//and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		//we create the form data field 'fileupload'
		//wich returns another writer to write the actual file
		partWriter, err := writer.CreateFormFile("uploadFile", "image.png")
		if err != nil {
			t.Error(err)
		}

		//https://yourbasic.org/golang/create-image/
		img := createImage()

		//Encode() takes an io.Writer.
		//We pass the multipart field
		//'fileupload' that we defined
		//earlier which, in turn, writes
		//to our io.Pipe
		err = png.Encode(partWriter, img)
		if err != nil {
			t.Error(err)
		}
	}()

	//We read from the pipe which receives data
	//from the multipart writer, which, in turn,
	//receives data from png.Encode().
	//We have 3 chained writers !
	request := httptest.NewRequest("POST", "/upload", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	response := recordRequest(request)

	t.Log("It should respond with an HTTP status code of 200")
	if response.Code != 200 {
		t.Errorf("Expected %d, received %d", 200, response.Code)
	}
	//t.Log("It should create a file named 'someimg.png' in uploads folder")
	//if _, err := os.Stat("./uploads/someimg.png"); os.IsNotExist(err) {
	//	t.Error("Expected file ./uploads/image.png' to exist")
	//}
}

func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{X: 0, Y: 0 }
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{R: 100, G: 200, B: 200, A: 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	_ = png.Encode(f, img)
	return img
}
