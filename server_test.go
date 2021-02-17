package crane_test

import (
	"context"
	"testing"

	"github.com/antklim/crane"
	"github.com/stretchr/testify/assert"
)

func TestCraneServer(t *testing.T) {
	var e crane.Event
	var err error

	ctx := context.Background()
	e = crane.Event{Category: "change", Commit: "abcd1234"}
	server := crane.New(nil)
	f := func() { server.Serve(ctx, e) }

	assert.PanicsWithValue(t, "crane: nil server handler", f)

	mux := crane.NewEventMux()
	mux.HandleFunc("change", func(c context.Context, e crane.Event) error { return nil })

	server.SetHandler(mux)
	assert.PanicsWithValue(t, "crane: multiple server handler registrations", func() {
		server.SetHandler(mux)
	})

	err = server.Serve(ctx, e)
	assert.NoError(t, err)

	e = crane.Event{Category: "change2", Commit: "abcd1234"}
	err = server.Serve(ctx, e)
	assert.EqualError(t, err, "crane: event handler not found")
}

func TestCraneEventMux(t *testing.T) {
	mux := crane.NewEventMux()
	assert.PanicsWithValue(t, "crane: invalid event category", func() { mux.Handle("", nil) })
	assert.PanicsWithValue(t, "crane: nil handler", func() { mux.Handle("change", nil) })
	assert.PanicsWithValue(t, "crane: nil handler", func() { mux.HandleFunc("change", nil) })
	mux.HandleFunc("change", func(c context.Context, e crane.Event) error { return nil })
	assert.PanicsWithValue(t, "crane: multiple registrations for change", func() {
		mux.HandleFunc("change", func(c context.Context, e crane.Event) error { return nil })
	})
}
