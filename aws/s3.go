package aws

import (
	"context"
	"path"

	"github.com/antklim/crane"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type bucketService struct {
	svc s3iface.S3API
}

// NewBucketClient creates a new instance of the bucket client.
func NewBucketClient(svc s3iface.S3API) crane.BucketAPI {
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

// CopyBucketWithContext copies objects of scrBucket/srcKey to destBucket/destKeyPrefix.
//
// When destKeyPrefix is "" then all source objects will be copied to the root
// of the destination bucket.
func (b *bucketService) CopyBucketWithContext(
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

			if _, err := b.svc.CopyObject(copyInput); err != nil {
				iterErr = err
				return false
			}
		}

		return !lastPage
	}

	// err is a pagination error
	if err := b.svc.ListObjectsV2Pages(input, iter); err != nil {
		return err
	}
	return iterErr
}
