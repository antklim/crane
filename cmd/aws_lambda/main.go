package main

import (
	"context"
	"log"

	"github.com/antklim/crane"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

var (
	version     string // sha1 of the code commit
	craneServer *crane.Server
)

func init() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	h := craneHandler(cfg)
	craneServer = crane.New(h)
}

func lambdaHandler(ctx context.Context, event crane.Event) error {
	log.Printf("crane: version %s\n", version)
	log.Printf("crane: event %+v\n", event)

	return craneServer.Serve(ctx, event)
}

func main() {
	runtime.Start(lambdaHandler)
}
