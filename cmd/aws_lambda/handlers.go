package main

import (
	"context"
	"log"
	"path"

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

		assetsSize, err := bc.KeySizeWithContext(ctx, cfg.DeployBucket, event.Commit)
		if err != nil {
			return errors.Wrap(err, "list changed files failed")
		}

		log.Printf("found %d of changed files", assetsSize)
		if assetsSize == 0 {
			return nil
		}

		archiveKeyPrefix := path.Join(cfg.ArchiveFolder, "pre_"+event.Commit)
		err = bc.CopyObjectsWithContext(ctx, cfg.StageBucket, "", cfg.ArchiveBucket, archiveKeyPrefix)
		if err != nil {
			return errors.Wrap(err, "archive stage bucket failed")
		}

		err = bc.SyncObjectsWithContext(ctx, cfg.DeployBucket, event.Commit, cfg.StageBucket, "")
		if err != nil {
			return errors.Wrap(err, "sync of assets and stage buckets failed")
		}

		return nil
	}
}

func releaseHandler(cfg *config, sess *session.Session) handlerFunc {
	return func(ctx context.Context, event crane.Event) error {
		bc := aws.NewBucketClient(s3.New(sess))

		err := bc.SyncObjectsWithContext(ctx, cfg.StageBucket, "", cfg.ProductionBucket, "")
		if err != nil {
			return errors.Wrap(err, "sync of assets and stage buckets failed")
		}

		return nil
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
