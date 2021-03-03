package main

import (
	"context"
	"testing"

	"github.com/antklim/crane"
	"github.com/stretchr/testify/assert"
)

func TestLambdaHandler(t *testing.T) {
	event := crane.Event{
		Category: "unknown",
		Commit:   "667e4625",
	}

	err := lambdaHandler(context.Background(), event)
	assert.EqualError(t, err, "crane: event handler not found")
}
