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

type handler func(*config) func(context.Context, crane.Event) error

var handlersMap = map[string]handler{
	changeCategory:   changeHandler,
	releaseCategory:  releaseHandler,
	rollbackCategory: rollbackHandler,
}

func changeHandler(cfg *config) func(ctx context.Context, event crane.Event) error {
	return func(ctx context.Context, event crane.Event) error {
		return errors.New("not implemented")
	}
}

func releaseHandler(cfg *config) func(ctx context.Context, event crane.Event) error {
	return func(ctx context.Context, event crane.Event) error {
		return errors.New("not implemented")
	}
}

func rollbackHandler(cfg *config) func(ctx context.Context, event crane.Event) error {
	return func(ctx context.Context, event crane.Event) error {
		return errors.New("not implemented")
	}
}

func craneHandler(cfg *config) *crane.EventMux {
	mux := crane.NewEventMux()

	for category, h := range handlersMap {
		mux.HandleFunc(category, h(cfg))
	}

	return mux
}
