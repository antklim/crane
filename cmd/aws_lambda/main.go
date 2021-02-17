package main

import (
	"context"

	"github.com/antklim/crane"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

var craneServer *crane.Server

func init() {
	h := craneHandler()
	craneServer = crane.New(h)
}

func handler(ctx context.Context, event crane.Event) error {
	// TODO: Logging

	return craneServer.Serve(ctx, event)
}

func main() {
	runtime.Start(handler)
}
