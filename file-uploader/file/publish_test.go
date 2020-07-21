package file

import (
	"context"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func CreateTestPubsub() pubsub.Client {
	ctx := context.Background()
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	defer conn.Close()
	// Use the connection when creating a pubsub client.
	client, _ := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	defer client.Close()
	return *client
}