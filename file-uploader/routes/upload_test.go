package routes

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/file"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/inject"
)

//func TestFileUpload(t *testing.T) {
//	processor := &FileProcessor{
//		Config: config.Config{},
//	}
//
//	pr, pw := io.Pipe()
//	//This writers is going to transform
//	//what we pass to it to multipart form data
//	//and write it to our io.Pipe
//	writer := multipart.NewWriter(pw)
//
//	go func() {
//		defer writer.Close()
//		//we create the form data field 'fileupload'
//		//wich returns another writer to write the actual file
//		part, err := writer.CreateFormFile("file", "sample_test_file.csv")
//		if err != nil {
//			t.Error(err)
//		}
//
//		//https://yourbasic.org/golang/create-image/
//		img := createImage()
//
//		//Encode() takes an io.Writer.
//		//We pass the multipart field
//		//'fileupload' that we defined
//		//earlier which, in turn, writes
//		//to our io.Pipe
//		err = png.Encode(part, img)
//		if err != nil {
//			t.Error(err)
//		}
//	}()
//}

var fileProcessorStub file.FileProcessor
var ctx = context.Background()
func init() {
	_, client := createTestPubSubServer("testtopic")
	fileProcessorStub = file.FileProcessor{
		Config: config.Config{
			Port: "8080",
			Pubsub: config.Pubsub{
				TopicId: "testtopic",
				ProjectId: "project",
			},
		},
		Client: client,
		Ctx: ctx,
	}
}


func createTestPubSubServer(topicId string) (*grpc.ClientConn, *pubsub.Client) {
	// Start a fake server running locally.
	srv := pstest.NewServer()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	// Use the connection when creating a pubsub client.
	client, _ := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	topic, _ := client.CreateTopic(ctx, topicId)
	_ = topic
	return conn, client
}

func TestFileUploadSuccess(t *testing.T) {
	inject.FileProcessor = fileProcessorStub
	path := "../file/sample_test_file.csv"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest("POST", "/samples/B/fileupload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	ProcessFile(res, req)

	assert.Equal(t, 202, res.Code)
}

//func TestFileUploadEmptyFileFail(t *testing.T) {
//	path := "./empty_csv.csv"
//	file, err := os.Open(path)
//	if err != nil {
//		t.Error(err)
//	}
//
//	defer file.Close()
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//	part, err := writer.CreateFormFile("file", filepath.Base(path))
//	if err != nil {
//		writer.Close()
//		t.Error(err)
//	}
//	io.Copy(part, file)
//	writer.Close()
//
//	req := httptest.NewRequest("POST", "/samples/B/fileupload", body)
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//	res := httptest.NewRecorder()
//
//	ProcessFile(res, req)
//
//	assert.Equal(t, 202, res.Code)
//}