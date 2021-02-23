package aws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antklim/crane/aws"
	"github.com/antklim/crane/aws/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKeySizeWithContext(t *testing.T) {
	s3ApiMock := mocks.S3API{}
	s3ApiMock.
		On("ListObjectsV2Pages", mock.AnythingOfType("*s3.ListObjectsV2Input"), mock.Anything).
		Return(errors.New("Test S3 API error"))
	bc := aws.NewBucketClient(&s3ApiMock)
	keySize, err := bc.KeySizeWithContext(context.Background(), "bucket", "key")
	assert.EqualError(t, err, "Test S3 API error")
	assert.Equal(t, 0, keySize)
}

func TestCopyObjectsWithContext(t *testing.T) {
	s3ApiMock := mocks.S3API{}
	s3ApiMock.
		On("ListObjectsV2Pages", mock.AnythingOfType("*s3.ListObjectsV2Input"), mock.Anything).
		Return(errors.New("Test S3 API error"))
	bc := aws.NewBucketClient(&s3ApiMock)
	err := bc.CopyObjectsWithContext(context.Background(), "srcBucket", "srcKey", "destBucket", "destKeyPrefix")
	assert.EqualError(t, err, "Test S3 API error")
}
