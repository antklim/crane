package crane

import (
	"context"
	"errors"
	"sync"
)

// Handler responds to event.
type Handler interface {
	Do(context.Context, Event) error
}

// HandlerFunc type is an adapter to allow the use of
// ordinary functions as event handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(context.Context, Event) error

// Do calls f(ctx, e).
func (f HandlerFunc) Do(ctx context.Context, e Event) error {
	return f(ctx, e)
}

var errNotFound = errors.New("crane: event handler not found")

// NotFound replies to the event with not found error.
func NotFound(ctx context.Context, e Event) error { return errNotFound }

// NotFoundHandler returns a simple event handler
// that replies to each event with not found error reply.
func NotFoundHandler() Handler { return HandlerFunc(NotFound) }

// EventMux is a crane event multiplexer.
// It matches the category of each event against a list of registered
// categories and calls the handler for the category.
type EventMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	h        Handler
	category string
}

// NewEventMux allocates and returns a new EventMux.
func NewEventMux() *EventMux { return new(EventMux) }

// Handler returns the handler to use for the given event,
// consulting event category.
//
// If there is no registered handler that applies to the event category,
// Handler returns a "category not found" handler.
func (mux *EventMux) Handler(e Event) Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	var h Handler
	if entry, ok := mux.m[e.Category]; ok {
		h = entry.h
	}
	if h == nil {
		h = NotFoundHandler()
	}
	return h
}

// Do dispatches the event to the caegory handler.
func (mux *EventMux) Do(ctx context.Context, e Event) error {
	h := mux.Handler(e)
	return h.Do(ctx, e)
}

// Handle registers the handler for the given category.
// If a handler already exists for category, Handle panics.
func (mux *EventMux) Handle(category string, h Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if category == "" {
		panic("crane: invalid event category")
	}
	if h == nil {
		panic("crane: nil handler")
	}
	if _, exists := mux.m[category]; exists {
		panic("crane: multiple registrations for " + category)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	mux.m[category] = muxEntry{h: h, category: category}
}

// HandleFunc registers the handler function for the given category.
func (mux *EventMux) HandleFunc(category string, h func(context.Context, Event) error) {
	if h == nil {
		panic("crane: nil handler")
	}
	mux.Handle(category, HandlerFunc(h))
}

// Server defines parameters for running a crane event server.
type Server struct {
	mu sync.Mutex
	h  Handler
}

// New creates a new Server.
func New(h Handler) *Server {
	return &Server{h: h}
}

// Serve processes an incoming event using registered handler.
func (s *Server) Serve(ctx context.Context, e Event) error {
	if s.h == nil {
		panic("crane: nil server handler")
	}

	return s.h.Do(ctx, e)
}

// SetHandler sets the server handler.
func (s *Server) SetHandler(h Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.h != nil {
		panic("crane: multiple server handler registrations")
	}

	s.h = h
}
