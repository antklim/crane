package main

//go:generate go run github.com/vektra/mockery/v2/ --name BucketAPI --srcpkg github.com/antklim/crane

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

type service struct {
	bc  crane.BucketAPI
	cfg *config
}

func newService(bc crane.BucketAPI, cfg *config) *service {
	return &service{bc: bc, cfg: cfg}
}

func (s *service) change(ctx context.Context, event crane.Event) error {
	assetsSize, err := s.bc.KeySizeWithContext(ctx, s.cfg.DeployBucket, event.Commit)
	if err != nil {
		return errors.Wrap(err, "list changed files failed")
	}

	log.Printf("found %d of changed files", assetsSize)
	if assetsSize == 0 {
		return nil
	}

	archiveKeyPrefix := path.Join(s.cfg.ArchiveFolder, "pre_"+event.Commit)
	err = s.bc.CopyObjectsWithContext(ctx, s.cfg.StageBucket, "", s.cfg.ArchiveBucket, archiveKeyPrefix)
	if err != nil {
		return errors.Wrap(err, "archive stage bucket failed")
	}

	err = s.bc.SyncObjectsWithContext(ctx, s.cfg.DeployBucket, event.Commit, s.cfg.StageBucket, "")
	if err != nil {
		return errors.Wrap(err, "sync of assets and stage buckets failed")
	}

	return nil
}

func (s *service) release(ctx context.Context) error {
	err := s.bc.SyncObjectsWithContext(ctx, s.cfg.StageBucket, "", s.cfg.ProductionBucket, "")
	if err != nil {
		return errors.Wrap(err, "sync of assets and stage buckets failed")
	}

	return nil
}

func changeHandler(cfg *config, sess *session.Session) handlerFunc {
	return func(ctx context.Context, event crane.Event) error {
		bc := aws.NewBucketClient(s3.New(sess))
		s := newService(bc, cfg)
		return s.change(ctx, event)
	}
}

func releaseHandler(cfg *config, sess *session.Session) handlerFunc {
	return func(ctx context.Context, event crane.Event) error {
		bc := aws.NewBucketClient(s3.New(sess))
		s := newService(bc, cfg)
		return s.release(ctx)
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
