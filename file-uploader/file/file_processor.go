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
	Client pubsub.Client
}

func (f *FileProcessor) ChunkCsv(file multipart.File, handler *multipart.FileHeader) {
	log.WithField("filename", handler.Filename).
		WithField("filesize", handler.Size).
		WithField("MIMEHeader", handler.Header).
		Info("File uploaded")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		log.WithField("line", line).
			Debug("Publishing csv line")
		err := f.Publish(line)
		if err != nil {
			log.WithField("line", line).
				Fatal("Error publishing csv line")
		}
	}

	if err := scanner.Err(); err != nil {
		log.WithError(err).
			Fatal("Error scanning file")
	}
}

func (f *FileProcessor) Publish(line string) error {
	topic := f.Client.Topic(f.Config.Pubsub.TopicId)

	ctx := context.Background()
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(line),
	})

	errorChannel := make(chan error, 1)
	go func(res *pubsub.PublishResult) {
		// The Get method blocks until a server-generated ID or
		// an error is returned for the published message.
		id, err := res.Get(ctx)
		if err != nil {
			// Error handling code can be added here.
			log.WithError(err).
				Fatal("Failed to publish")
			errorChannel <- err
		}
		log.WithField("line", line).
			WithField("messageId", id).
			Debug("published message")
		errorChannel <- nil
	}(result)

	return <- errorChannel
}
