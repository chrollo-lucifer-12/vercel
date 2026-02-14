package storage

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chrollo-lucifer-12/shared/utils"
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

func (s *S3Storage) UploadDirectory(ctx context.Context, localDir, slug string) error {

	baseDir, err := filepath.Abs(localDir)
	if err != nil {
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
		return err
	}

	log.Printf("Found %d files to upload\n", len(files))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, filePath := range files {
		filePath := filePath

		g.Go(func() error {
			return s.uploadSingleFile(ctx, baseDir, filePath, slug)
		})
	}

	return g.Wait()
}

func (s *S3Storage) uploadSingleFile(ctx context.Context, baseDir, filePath, slug string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return err
	}

	objectKey := filepath.ToSlash(filepath.Join(slug, relPath))
	contentType := utils.DetectContentType(objectKey)

	log.Printf("Uploading %s → %s\n", filePath, objectKey)

	size := stat.Size()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentType:   aws.String(contentType),
		ContentLength: &size,
	})

	return err
}
