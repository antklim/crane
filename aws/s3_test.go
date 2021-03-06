package aws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antklim/crane/aws"
	"github.com/antklim/crane/aws/mocks"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKeySizeWithContext(t *testing.T) {
	bucket := "bucket"
	key := "key"

	s3ApiMock := mocks.S3API{}
	s3ApiMock.On(
		"ListObjectsV2Pages",
		&s3.ListObjectsV2Input{
			Bucket: &bucket,
			Prefix: &key,
		},
		mock.Anything,
	).Return(errors.New("Test S3 API error")).Once()

	bc := aws.NewBucketClient(&s3ApiMock)
	keySize, err := bc.KeySizeWithContext(context.Background(), bucket, key)
	assert.EqualError(t, err, "Test S3 API error")
	assert.Equal(t, 0, keySize)
	s3ApiMock.AssertExpectations(t)
}

func TestCopyObjectsWithContext(t *testing.T) {
	srcBucket := "srcBucket"
	srcKey := "srcKey"
	destBucket := "destBucket"
	destKeyPrefix := "destKeyPrefix"

	s3ApiMock := mocks.S3API{}
	s3ApiMock.On(
		"ListObjectsV2Pages",
		&s3.ListObjectsV2Input{
			Bucket: &srcBucket,
			Prefix: &srcKey,
		},
		mock.Anything,
	).Return(errors.New("Test S3 API error")).Once()

	bc := aws.NewBucketClient(&s3ApiMock)
	err := bc.CopyObjectsWithContext(context.Background(), srcBucket, srcKey, destBucket, destKeyPrefix)
	assert.EqualError(t, err, "Test S3 API error")
	s3ApiMock.AssertExpectations(t)
}
