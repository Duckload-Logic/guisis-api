package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LightsailStorage struct {
	s3Client   *s3.Client
	bucketName string
}

func NewLightsailStorage(
	ctx context.Context,
	bucketName string,
	region string,
) (*LightsailStorage, error) {
	if bucketName == "" || region == "" {
		return nil, fmt.Errorf(
			"[NewLightsailStorage] {Initialize}: " +
				"bucket name and region must be defined",
		)
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[NewLightsailStorage] {Initialize}: " +
				"unable to load AWS SDK config: %w",
			err,
		)
	}

	s3Client := s3.NewFromConfig(cfg)

	return &LightsailStorage{
		s3Client:   s3Client,
		bucketName: bucketName,
	}, nil
}

// Upload uploads data from a reader to the given blob path.
func (l *LightsailStorage) Upload(
	ctx context.Context,
	path string,
	reader io.ReadSeeker,
	contentType string,
) error {
	_, err := l.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(l.bucketName),
		Key:         aws.String(path),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf(
			"[LightsailStorage] {Upload}: "+
				"failed to upload blob %q: %w",
			path,
			err,
		)
	}
	return nil
}

// Download streams the blob content into the provided writer.
func (l *LightsailStorage) Download(
	ctx context.Context,
	path string,
	writer io.Writer,
) error {
	resp, err := l.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(l.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf(
			"[LightsailStorage] {Download}: "+
				"failed to download blob %q: %w",
			path,
			err,
		)
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf(
			"[LightsailStorage] {Download}: "+
				"failed to copy blob data %q: %w",
			path,
			err,
		)
	}
	return nil
}

// Delete removes the blob from the container.
func (l *LightsailStorage) Delete(ctx context.Context, path string) error {
	_, err := l.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(l.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf(
			"[LightsailStorage] {Delete}: "+
				"failed to delete blob %q: %w",
			path,
			err,
		)
	}
	return nil
}
