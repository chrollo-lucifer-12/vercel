package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chrollo-lucifer-12/build-server/logs"
	"github.com/chrollo-lucifer-12/shared/utils"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(endpoint, accessKey, secretKey, region, bucket string) (*S3Storage, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &S3Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Storage) UploadDirectory(ctx context.Context, localDir, slug string, dispatcher *logs.LogDispatcher, deploymentIDUUID uuid.UUID) error {
	dispatcher.Push(deploymentIDUUID, "Starting directory upload to S3...")
	baseDir, err := filepath.Abs(localDir)
	if err != nil {
		dispatcher.Push(deploymentIDUUID, "Failed resolving directory path: "+err.Error())
		return err
	}

	var files []string

	err = filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		dispatcher.Push(deploymentIDUUID, "Failed walking directory: "+err.Error())
		return err
	}
	dispatcher.Push(deploymentIDUUID,
		fmt.Sprintf("Found %d files to upload...", len(files)),
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, filePath := range files {
		filePath := filePath

		g.Go(func() error {
			err := s.uploadSingleFile(ctx, baseDir, filePath, slug, dispatcher, deploymentIDUUID)
			if err != nil {
				dispatcher.Push(deploymentIDUUID, "Upload failed: "+filePath+" -> "+err.Error())
			}
			return err
		})
	}

	return g.Wait()
}

func (s *S3Storage) uploadSingleFile(
	ctx context.Context,
	baseDir, filePath, slug string,
	dispatcher *logs.LogDispatcher,
	deploymentIDUUID uuid.UUID,
) error {

	relPath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return err
	}

	objectKey := filepath.ToSlash(filepath.Join(slug, relPath))

	dispatcher.Push(deploymentIDUUID, "Uploading: "+objectKey)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	contentType := utils.DetectContentType(objectKey)
	size := stat.Size()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentType:   aws.String(contentType),
		ContentLength: &size,
	})

	if err != nil {
		return err
	}

	dispatcher.Push(deploymentIDUUID, "Uploaded successfully: "+objectKey)
	return nil
}

func (s *S3Storage) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {

	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return out.Body, nil
}
