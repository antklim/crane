package crane

import "context"

// BucketAPI defines generic bucket operations.
type BucketAPI interface {
	KeySizeWithContext(ctx context.Context, bucket, key string) (int, error)
	CopyObjectsWithContext(ctx context.Context, srcBucket, srcKey, destBucket, destKeyPrefix string) error
	DeleteObjectsWithContext(ctx context.Context, bucket, keyPrefix string) error
	SyncObjectsWithContext(ctx context.Context, srcBucket, srcKey, destBucket, destKeyPrefix string) error
}
