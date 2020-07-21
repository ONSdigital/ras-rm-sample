package file

import (
	"bufio"
	"cloud.google.com/go/pubsub"
	"context"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
)

type FileProcessor struct {
	Config config.Config
	Client *pubsub.Client
	Ctx context.Context
}

func (f *FileProcessor) ChunkCsv(file multipart.File, handler *multipart.FileHeader) {
	log.WithField("filename", handler.Filename).
		WithField("filesize", handler.Size).
		WithField("MIMEHeader", handler.Header).
		Info("File uploaded")
	f.Publish(bufio.NewScanner(file))
}

func (f *FileProcessor) Publish(scanner *bufio.Scanner) int {
	topic := f.Client.Topic(f.Config.Pubsub.TopicId)
	var errorCount = 0
	for scanner.Scan() {
		line := scanner.Text()
		log.WithField("line", line).
			Debug("Publishing csv line")

		errorChannel := make(chan error, 0)

		go func(line string, topic *pubsub.Topic) {
			id, err := topic.Publish(f.Ctx, &pubsub.Message{
				Data: []byte(line),
			}).Get(f.Ctx)
			log.WithField("line", line).
				WithField("messageId", id).
				Debug("csv line acknowledged")
			errorChannel <- err
		}(line, topic)

		err := <- errorChannel
		if err != nil {
			errorCount++
			log.WithField("line", line).
				WithError(err).
				Error("Error publishing csv line")
		}
	}
	if err := scanner.Err(); err != nil {
		log.WithError(err).Error("Error scanning file")
	}
	return errorCount
}
