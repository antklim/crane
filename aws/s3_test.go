package aws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antklim/crane/aws"
	"github.com/antklim/crane/aws/mocks"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKeySizeWithContext(t *testing.T) {
	testCases := []struct {
		desc   string
		s3api  func() s3iface.S3API
		assert func(*testing.T, int, error)
	}{
		{
			desc: "propagates S3 API errors",
			s3api: func() s3iface.S3API {
				apiMock := mocks.S3API{}
				apiMock.
					On("ListObjectsV2Pages", mock.AnythingOfType("*s3.ListObjectsV2Input"), mock.Anything).
					Return(errors.New("Test S3 API error"))
				return &apiMock
			},
			assert: func(t *testing.T, keySize int, err error) {
				assert.EqualError(t, err, "Test S3 API error")
				assert.Equal(t, 0, keySize)
			},
		},
		{
			desc: "returns bucket key objects amount",
			s3api: func() s3iface.S3API {
				apiMock := mocks.S3API{}
				apiMock.
					On("ListObjectsV2Pages", mock.AnythingOfType("*s3.ListObjectsV2Input"), mock.Anything).
					Return(nil)
				return &apiMock
			},
			assert: func(t *testing.T, keySize int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 0, keySize)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bc := aws.NewBucketClient(tC.s3api())
			keySize, err := bc.KeySizeWithContext(context.Background(), "bucket", "key")
			tC.assert(t, keySize, err)
		})
	}
}

func TestCopyObjectsWithContext(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}

func TestDeleteObjectsWithContext(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}

func TestSyncObjectsWithContext(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
