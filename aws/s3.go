package aws

//go:generate go run github.com/vektra/mockery/v2/ --name S3API --srcpkg github.com/aws/aws-sdk-go/service/s3/s3iface

import (
	"context"
	"path"

	"github.com/antklim/crane"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// TODO: Add integration tests to validate interator functions.
// TODO: Add support of copy bucket objects with subdirectories.
// Currently, copy flattens the structure.

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
	err := s.svc.ListObjectsV2Pages(input, keySizeIterator(&keySize))

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
	iter := copyObjectsIterator(s.svc, srcBucket, destBucket, destKeyPrefix, &iterErr)

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

	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(keyPrefix),
	}

	iter := s3manager.NewDeleteListIterator(s.svc, input)

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

func keySizeIterator(counter *int) func(*s3.ListObjectsV2Output, bool) bool {
	return func(out *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, v := range out.Contents {
			// out.Contents contains directories keys that have no size
			if aws.Int64Value(v.Size) > 0 {
				*counter++
			}
		}
		return !lastPage
	}
}

func copyObjectsIterator(
	svc s3iface.S3API,
	srcBucket, destBucket, destKeyPrefix string,
	iterErr *error,
) func(*s3.ListObjectsV2Output, bool) bool {

	return func(out *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, v := range out.Contents {
			if aws.Int64Value(v.Size) == 0 {
				continue
			}

			destKey := path.Join(destKeyPrefix, path.Base(*v.Key))
			src := path.Join(srcBucket, *v.Key)

			input := &s3.CopyObjectInput{
				Bucket:     aws.String(destBucket),
				Key:        aws.String(destKey),
				CopySource: aws.String(src),
			}

			if _, err := svc.CopyObject(input); err != nil {
				*iterErr = err
				return false
			}
		}

		return !lastPage
	}
}
