package main

import (
	"context"
	"log"

	"github.com/antklim/crane"
	"github.com/antklim/crane/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

const (
	changeCategory   string = "change"
	releaseCategory  string = "release"
	rollbackCategory string = "rollback"
)

type handlerFunc func(context.Context, crane.Event) error
type handler func(*config, *session.Session) handlerFunc

var handlersMap = map[string]handler{
	changeCategory:   changeHandler,
	releaseCategory:  releaseHandler,
	rollbackCategory: rollbackHandler,
}

func changeHandler(cfg *config, sess *session.Session) handlerFunc {
	return func(ctx context.Context, event crane.Event) error {
		bc := aws.NewBucketClient(s3.New(sess))

		assetsSize, err := bc.BucketKeySizeWithContext(ctx, cfg.DeployBucket, event.Commit)
		if err != nil {
			return errors.Wrap(err, "list changed files failed")
		}

		log.Printf("found %d of changed files", assetsSize)

		if assetsSize == 0 {
			return nil
		}

		return nil
	}
}

func releaseHandler(cfg *config, sess *session.Session) handlerFunc {
	return func(ctx context.Context, event crane.Event) error {
		return errors.New("not implemented")
	}
}

func rollbackHandler(cfg *config, sess *session.Session) handlerFunc {
	return func(ctx context.Context, event crane.Event) error {
		return errors.New("not implemented")
	}
}

func craneHandler(cfg *config, sess *session.Session) *crane.EventMux {
	mux := crane.NewEventMux()

	for category, h := range handlersMap {
		mux.HandleFunc(category, h(cfg, sess))
	}

	return mux
}
