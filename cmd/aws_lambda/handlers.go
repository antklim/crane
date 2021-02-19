package main

import (
	"context"
	"errors"
	"log"

	"github.com/antklim/crane"
	"github.com/aws/aws-sdk-go/aws/session"
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
		log.Println("Listing assets...")
		// svc := s3.New(sess)
		// assets, err := listDeployAssets(svc, cfg.DeployBucket, event.Commit)
		// if err != nil {
		// 	log.Println("Assets list failed")
		// 	return err
		// }

		// if len(assets) == 0 {
		// 	log.Println("No assets found for deployment")
		// 	return nil
		// }

		// log.Printf("Found %d assets\n", len(assets))

		log.Println("Archiving target...")
		// if err := archiveTarget(svc, cfg.StageBucket, cfg.ArchiveBucket, cfg.ArchiveFolder, event.Commit); err != nil {
		// 	log.Println("Target archive failed")
		// 	return err
		// }
		log.Println("Target archived")

		log.Println("Syncing target...")
		// if err := syncTarget(svc, cfg.StageBucket, cfg.DeployBucket, event.Commit); err != nil {
		// 	log.Println("Target sync failed")
		// 	return err
		// }

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
