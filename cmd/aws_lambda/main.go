package main

import (
	"context"
	"log"

	"github.com/antklim/crane"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	version     string // sha1 of the code commit
	cfg         *config
	sess        *session.Session
	craneServer *crane.Server
)

func init() {
	var err error

	cfg, err = loadConfig()
	if err != nil {
		log.Panic(err)
	}

	sess, err = session.NewSession(&aws.Config{Region: aws.String(cfg.Region)})
	if err != nil {
		log.Panic(err)
	}

	h := craneHandler(cfg, sess)
	craneServer = crane.New(h)
}

func lambdaHandler(ctx context.Context, event crane.Event) error {
	log.Printf("version %s\n", version)
	log.Printf("event %+v\n", event)

	return craneServer.Serve(ctx, event)
}

func main() {
	runtime.Start(lambdaHandler)
}
