package aws

import (
	"errors"
	"testing"

	"github.com/antklim/crane/aws/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

func TestKeySizeIterator(t *testing.T) {
	counter := 0
	p1out := &s3.ListObjectsV2Output{
		Contents: []*s3.Object{
			{
				Key:          aws.String("3342d1f9/"),
				ETag:         aws.String("\"d41d8cd98f00b204e9800998ecf8427e\""),
				Size:         aws.Int64(0),
				StorageClass: aws.String("STANDARD"),
			},
			{
				Key:          aws.String("3342d1f9/test.txt"),
				ETag:         aws.String("\"a135c23302404725d97783611628f077\""),
				Size:         aws.Int64(12),
				StorageClass: aws.String("STANDARD"),
			},
		},
	}

	p2out := &s3.ListObjectsV2Output{
		Contents: []*s3.Object{
			{
				Key:          aws.String("3342d1f9/foo/"),
				ETag:         aws.String("\"d41d8cd98f00b204e9800998ecf8427e\""),
				Size:         aws.Int64(0),
				StorageClass: aws.String("STANDARD"),
			},
			{
				Key:          aws.String("3342d1f9/foo/event.go"),
				ETag:         aws.String("\"35da3ec4cc62b55e9bc5cc083c1dbce0\""),
				Size:         aws.Int64(230),
				StorageClass: aws.String("STANDARD"),
			},
		},
	}

	lastPage := false
	shouldContinue := keySizeIterator(&counter)(p1out, lastPage)
	assert.True(t, shouldContinue)
	assert.Equal(t, 1, counter)

	lastPage = true
	shouldContinue = keySizeIterator(&counter)(p2out, lastPage)
	assert.False(t, shouldContinue)
	assert.Equal(t, 2, counter)
}

func TestCopyObjectsIterator(t *testing.T) {
	srcBucket := "srcBucket"
	destBucket := "destBucket"
	destKeyPrefix := "destKeyPrefix"

	destKey := "destKeyPrefix/test.txt"
	src := "srcBucket/3342d1f9/test.txt"

	out := &s3.ListObjectsV2Output{
		Contents: []*s3.Object{
			{
				Key:          aws.String("3342d1f9/"),
				ETag:         aws.String("\"d41d8cd98f00b204e9800998ecf8427e\""),
				Size:         aws.Int64(0),
				StorageClass: aws.String("STANDARD"),
			},
			{
				Key:          aws.String("3342d1f9/test.txt"),
				ETag:         aws.String("\"a135c23302404725d97783611628f077\""),
				Size:         aws.Int64(12),
				StorageClass: aws.String("STANDARD"),
			},
		},
	}

	testCases := []struct {
		desc     string
		svc      func() (s3iface.S3API, func(*testing.T))
		lastPage bool
		assert   func(*testing.T, bool, error)
	}{
		{
			desc: "stops iterator when failed to copy object",
			svc: func() (s3iface.S3API, func(*testing.T)) {
				m := &mocks.S3API{}
				m.On(
					"CopyObject",
					&s3.CopyObjectInput{
						Bucket:     &destBucket,
						Key:        &destKey,
						CopySource: &src,
					},
				).Return(nil, errors.New("CopyObject error")).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			lastPage: false,
			assert: func(t *testing.T, shouldContinue bool, err error) {
				assert.EqualError(t, err, "CopyObject error")
				assert.False(t, shouldContinue)
			},
		},
		{
			desc:     "stops iterator when reached the last page",
			lastPage: true,
			svc: func() (s3iface.S3API, func(*testing.T)) {
				m := &mocks.S3API{}
				m.On(
					"CopyObject",
					&s3.CopyObjectInput{
						Bucket:     &destBucket,
						Key:        &destKey,
						CopySource: &src,
					},
				).Return(nil, nil).Once()

				ae := func(t *testing.T) {
					m.AssertExpectations(t)
				}

				return m, ae
			},
			assert: func(t *testing.T, shouldContinue bool, err error) {
				assert.NoError(t, err)
				assert.False(t, shouldContinue)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			svc, ae := tC.svc()
			var err error
			shouldContinue := copyObjectsIterator(svc, srcBucket, destBucket, destKeyPrefix, &err)(out, tC.lastPage)
			tC.assert(t, shouldContinue, err)
			ae(t)
		})
	}
}
