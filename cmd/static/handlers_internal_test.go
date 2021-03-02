package main

import (
	"context"
	"testing"

	"github.com/antklim/crane"
	"github.com/antklim/crane/cmd/static/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: Add tests

func TestServiceChange(t *testing.T) {
	cfg := &config{
		ArchiveBucket:    "archive-bucket",
		ArchiveFolder:    "archive-folder",
		DeployBucket:     "deploy-bucket",
		StageBucket:      "stage-bucket",
		ProductionBucket: "production-bucket",
		Region:           "ap-southeast-1",
	}

	event := crane.Event{
		Category: changeCategory,
		Commit:   "667e4622",
	}

	testCases := []struct {
		desc   string
		bc     func() (crane.BucketAPI, func(*testing.T))
		assert func(*testing.T, error)
	}{
		{
			desc: "returns error when BucketAPI KeySizeWithContext failed",
			bc: func() (crane.BucketAPI, func(*testing.T)) {
				m := &mocks.BucketAPI{}
				m.On(
					"KeySizeWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
				).Return(0, errors.New("KeySizeWithContext error")).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "list changed files failed: KeySizeWithContext error")
			},
		},
		{
			desc: "returns nil when BucketAPI KeySizeWithContext returns 0",
			bc: func() (crane.BucketAPI, func(*testing.T)) {
				m := &mocks.BucketAPI{}
				m.On(
					"KeySizeWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
				).Return(0, nil).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			desc: "returns error when BucketAPI CopyObjectsWithContext failed",
			bc: func() (crane.BucketAPI, func(*testing.T)) {
				m := &mocks.BucketAPI{}
				m.On(
					"KeySizeWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
				).Return(1, nil).Once()
				m.On(
					"CopyObjectsWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"stage-bucket",
					"",
					"archive-bucket",
					"archive-folder/pre_667e4622",
				).Return(errors.New("CopyObjectsWithContext error")).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "archive stage bucket failed: CopyObjectsWithContext error")
			},
		},
		{
			desc: "returns error when BucketAPI SyncObjectsWithContext failed",
			bc: func() (crane.BucketAPI, func(*testing.T)) {
				m := &mocks.BucketAPI{}
				m.On(
					"KeySizeWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
				).Return(1, nil).Once()
				m.On(
					"CopyObjectsWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"stage-bucket",
					"",
					"archive-bucket",
					"archive-folder/pre_667e4622",
				).Return(nil).Once()
				m.On(
					"SyncObjectsWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
					"stage-bucket",
					"",
				).Return(errors.New("SyncObjectsWithContext error")).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "sync of assets and stage buckets failed: SyncObjectsWithContext error")
			},
		},
		{
			desc: "returns nil when successfully fnished",
			bc: func() (crane.BucketAPI, func(*testing.T)) {
				m := &mocks.BucketAPI{}
				m.On(
					"KeySizeWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
				).Return(1, nil).Once()
				m.On(
					"CopyObjectsWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"stage-bucket",
					"",
					"archive-bucket",
					"archive-folder/pre_667e4622",
				).Return(nil).Once()
				m.On(
					"SyncObjectsWithContext",
					mock.AnythingOfType("*context.emptyCtx"),
					"deploy-bucket",
					"667e4622",
					"stage-bucket",
					"",
				).Return(nil).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bc, ae := tC.bc()
			svc := newService(bc, cfg)
			err := svc.change(context.Background(), event)
			tC.assert(t, err)
			ae(t)
		})
	}
}

func TestServiceRelease(t *testing.T) {
	t.Skip("tests not implemented")
}
