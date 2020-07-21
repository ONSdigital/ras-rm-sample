package file

import (
	"context"
	//"fmt"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"testing"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var fileProcessorStub = &FileProcessor{
	Config: config.Config{
		Port: "8080",
		Pubsub: config.Pubsub{
			TopicId: "testtopic",
			ProjectId: "project",
	    },
	},
	Client: nil,
	Ctx: textContext,
}

var textContext = context.Background()

func TestSendingToPubsub(t *testing.T) {
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	defer conn.Close()
	// Use the connection when creating a pubsub client.
	client, _ := pubsub.NewClient(textContext, "project", option.WithGRPCConn(conn))
	topic, _ := client.CreateTopic(textContext, "testtopic")
	_ = topic
	defer client.Close()

	fileProcessorStub.Client = client

	err := fileProcessorStub.Publish("TestCSVLine")

	if err != nil {
		t.Errorf("Unexpected error thrown. expected: %v, actual: %v", nil, err)
	}
}

func TestSendToWrongTopicThrowsError(t *testing.T) {
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	defer conn.Close()
	// Use the connection when creating a pubsub client.
	client, _ := pubsub.NewClient(textContext, "project", option.WithGRPCConn(conn))
	topic, _ := client.CreateTopic(textContext, "badtopic")
	_ = topic
	defer client.Close()

	fileProcessorStub.Client = client

	err := fileProcessorStub.Publish("TestCSVLine")

	if err == nil {
		t.Errorf("Error expected. expected: %v, actual: %v", err, nil)
	}
}