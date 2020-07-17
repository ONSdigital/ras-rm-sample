package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type CSVWorker struct {
	sample []byte
}

func init() {
	// Only log the warning severity or above.
	verbose, err := strconv.ParseBool(getEnv("VERBOSE", "true"))
	if err != nil {
		log.Error("unable to parse verbose env")
		verbose = true
	}
	if verbose {
		//anything debug and above
		log.SetLevel(log.DebugLevel)
	} else {
		//otherwise keep it to info
		log.SetLevel(log.InfoLevel)
	}
}

func (cw CSVWorker) start() {
	log.Debug("starting worker process")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, getEnv("GOOGLE_CLOUD_PROJECT", "rm-ras-sandbox"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	log.Debug("about to subscribe")
	cw.subscribe(ctx, client)
}

func (cw CSVWorker) subscribe(ctx context.Context, client *pubsub.Client) {
	subId := getEnv("PUBSUB_SUB_ID", "sample-workers")
	log.WithField("subId", subId).Info("subscribing to subscription")
	sub := client.Subscription(subId)
	cctx, cancel := context.WithCancel(ctx)
	log.Debug("waiting to receive")
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Info("sample received - processing")
		log.WithField("data", string(msg.Data)).Debug("sample data")
		cw.sample = msg.Data
		err := processSample(cw.sample)
		if err != nil {
			log.WithError(err).Error("error processing sample - nacking message")
			msg.Nack()
		} else {
			log.Info("sample processed - acking message")
			msg.Ack()
		}
	})

	if err != nil {
		log.WithError(err).Error("error subscribing")
		cancel()
	}
}

func getEnv(key string, defaultVar string) string {
	v := os.Getenv(key)
	if v == "" {
		log.WithFields(log.Fields{
			"key":     key,
			"default": defaultVar,
		}).Info("environment variable not set using default")
		return defaultVar
	}
	return v
}

func main() {
	workers := getEnv("WORKERS", "10")
	noOfWorkers, err := strconv.Atoi(workers)
	if err != nil {
		log.Error("Unable to set number of workers defaulting to 10")
		noOfWorkers = 10
	}
	for i := 0; i < noOfWorkers; i++ {
		csvWorker := &CSVWorker{}
		go csvWorker.start()
	}
	select {
	// loop forever
	}
}
