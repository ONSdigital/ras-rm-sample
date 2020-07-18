package file

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go/scanner"
	"mime/multipart"
	"net/http"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
)

type FileProcessor struct {
	Config config.Config
}

func (f *FileProcessor) ChunkCsv(file multipart.File, handler *multipart.FileHeader) {
	log.WithField("filename", handler.Filename).
		WithField("filesize", handler.Size).
		WithField("MIMEHeader", handler.Header).
		Info("File uploaded")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.WithError(err).
			Fatal("Error scanning file")
	}
}

func ProcessFile(w http.ResponseWriter, r *http.Request) {
	// 10MB maximum file size
	//r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.WithError(err).
			Fatal("Error retrieving the file")
		return
	}
	processor := $FileProcessor{
		Config: &config.Config{},
	}
	processor.ChunkCsv(file, handler)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
