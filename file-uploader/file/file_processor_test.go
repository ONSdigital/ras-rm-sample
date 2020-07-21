package file

import (
	"bufio"
	"context"
	"log"
	"os"

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

func TestScannerAndPublishSuccess(t *testing.T) {
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

	file, err := os.Open("sample_test_file.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	errorCount := fileProcessorStub.Publish(scanner)

	if errorCount != 0 {
		t.Errorf("Errors have been thrown. expected: %v, actual: %v", 0, errorCount)
	}
}

func TestScannerAndPublishBadTopic(t *testing.T) {
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	defer conn.Close()
	// Use the connection when creating a pubsub client.
	client, _ := pubsub.NewClient(textContext, "project", option.WithGRPCConn(conn))
	topic, _ := client.CreateTopic(textContext, "BadTopictopic")
	_ = topic
	defer client.Close()

	fileProcessorStub.Client = client

	file, err := os.Open("sample_test_file.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	errorCount := fileProcessorStub.Publish(scanner)

	if errorCount != 8 {
		t.Errorf("Invalid amount of errors thrown. expected: %v, actual: %v", 8, errorCount)
	}
}