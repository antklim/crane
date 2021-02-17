package main

import (
	"context"
	"errors"

	"github.com/antklim/crane"
)

const (
	changeCategory   string = "change"
	releaseCategory  string = "release"
	rollbackCategory string = "rollback"
)

func changeHandler(ctx context.Context, event crane.Event) error {
	return errors.New("not implemented")
}

func releaseHandler(ctx context.Context, event crane.Event) error {
	return errors.New("not implemented")
}

func rollbackHandler(ctx context.Context, event crane.Event) error {
	return errors.New("not implemented")
}

func craneHandler() *crane.EventMux {
	mux := crane.NewEventMux()
	mux.HandleFunc(changeCategory, changeHandler)
	mux.HandleFunc(releaseCategory, releaseHandler)
	mux.HandleFunc(rollbackCategory, rollbackHandler)
	return mux
}
