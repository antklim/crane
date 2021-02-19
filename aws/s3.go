package aws

import (
	"context"
	"errors"

	"github.com/antklim/crane"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type bucketService struct {
	svc *s3.S3
}

// NewBucketClient creates a new instance of the bucket client.
func NewBucketClient(svc *s3.S3) crane.BucketAPI {
	return &bucketService{svc: svc}
}

// BucketKeySizeWithContext returns amount of objects nested in the key.
func (b *bucketService) BucketKeySizeWithContext(
	ctx context.Context, bucket, key string) (int, error) {

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}

	keySize := 0
	iter := func(out *s3.ListObjectsV2Output, lastPage bool) bool {
		keySize += len(out.Contents)
		return !lastPage
	}

	err := b.svc.ListObjectsV2Pages(input, iter)
	return keySize, err
}

func (b *bucketService) CopyBucketWithContext(
	ctx context.Context, srcBucket, srcKey, destBucket, destKey string) error {

	return errors.New("not implemented")
}
