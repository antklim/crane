package aws

import (
	"context"
	"path"

	"github.com/antklim/crane"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate go run github.com/vektra/mockery/v2/ --name S3API --srcpkg github.com/aws/aws-sdk-go/service/s3/s3iface

// TODO: Add integration tests to validate interator functions.

type bucketService struct {
	svc s3iface.S3API
}

// NewBucketClient creates a new instance of the bucket client.
func NewBucketClient(svc s3iface.S3API) crane.BucketAPI {
	return &bucketService{svc: svc}
}

// KeySizeWithContext returns amount of objects nested in the bucket key.
func (s *bucketService) KeySizeWithContext(
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

	err := s.svc.ListObjectsV2Pages(input, iter)
	return keySize, err
}

// CopyObjectsWithContext copies objects of scrBucket/srcKey to destBucket/destKeyPrefix.
//
// When destKeyPrefix is "" then all source objects will be copied to the root
// of the destination bucket.
func (s *bucketService) CopyObjectsWithContext(
	ctx context.Context, srcBucket, srcKey, destBucket, destKeyPrefix string) error {

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(srcBucket),
		Prefix: aws.String(srcKey), // TODO: test what happens when srcKey is ""
	}

	// TODO: use go-routines to copy objects
	// use errors channel to propagate errors from go routine

	var iterErr error
	iter := func(out *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, v := range out.Contents {
			destKey := path.Join(destKeyPrefix, path.Base(*v.Key))
			src := path.Join(srcBucket, *v.Key)

			copyInput := &s3.CopyObjectInput{
				Bucket:     aws.String(destBucket),
				Key:        aws.String(destKey),
				CopySource: aws.String(src),
			}

			if _, err := s.svc.CopyObject(copyInput); err != nil {
				iterErr = err
				return false
			}
		}

		return !lastPage
	}

	// err is a pagination error
	if err := s.svc.ListObjectsV2Pages(input, iter); err != nil {
		return err
	}
	return iterErr
}

// DeleteObjectsWithContext deletes all objects of bucket/keyPrefix.
//
// When keyPrefix is "" then all bucket objects will be deleted.
func (s *bucketService) DeleteObjectsWithContext(
	ctx context.Context, bucket, keyPrefix string) error {

	iter := s3manager.NewDeleteListIterator(s.svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(keyPrefix),
	})

	return s3manager.NewBatchDeleteWithClient(s.svc).Delete(ctx, iter)
}

// SyncObjectsWithContext overwrites all objects of destBucket/destKeyPrefix with
// the bjects of srcBucket/srcKey.
//
// When destKeyPrefix is "" then all source objects will be copied to the root
// of the destination bucket.
func (s *bucketService) SyncObjectsWithContext(
	ctx context.Context, srcBucket, srcKey, destBucket, destKeyPrefix string) error {

	if err := s.DeleteObjectsWithContext(ctx, destBucket, destKeyPrefix); err != nil {
		return err
	}

	return s.CopyObjectsWithContext(ctx, srcBucket, srcKey, destBucket, destKeyPrefix)
}
