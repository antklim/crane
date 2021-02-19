package crane

import "context"

// BucketAPI defines generic bucket operations.
type BucketAPI interface {
	BucketKeySizeWithContext(ctx context.Context, bucket, key string) (int, error)
	CopyBucketWithContext(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) error
}
