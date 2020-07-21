package file

import (
	//"fmt"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	log "github.com/sirupsen/logrus"
	"testing"
)

var fileProcessorStub = &FileProcessor{
	Config: config.Config{
		Port: "8080",
		Pubsub: config.Pubsub{
			TopicId: "test-topic",
			ProjectId: "project",
	    },
	},
	Client: CreateTestPubsub(),
}

func TestSendingToPubsub(t *testing.T) {
	err := fileProcessorStub.Publish("TestCSVLine")

	if err != nil {
		log.WithError(err).Fatal("err")
	}
}