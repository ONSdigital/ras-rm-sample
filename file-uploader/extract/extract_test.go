package extract

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
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

func TestFileUploadSuccess(t *testing.T) {
	path := "./sample_test_file.csv"
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

func TestFileUploadEmptyFileFail(t *testing.T) {
	path := "./empty_csv.csv"
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