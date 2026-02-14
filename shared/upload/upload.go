package upload

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/sync/errgroup"
)

type StorageInterface interface {
	UploadFile(bucketID, path string, reader io.Reader) error
	ListFiles(bucketID, path string) ([]string, error)
	RemoveFile(bucketID string, paths []string) error
}

type MinioStorage struct {
	Client *minio.Client
}

func NewMinioStorage(endpoint, accessKey, region, secretKey string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorage{Client: client}, nil
}

func (m *MinioStorage) UploadFile(bucketID, path string, reader io.Reader) error {
	ctx := context.Background()

	contentType := "application/octet-stream"
	if strings.HasSuffix(path, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript"
	}

	_, err := m.Client.PutObject(ctx, bucketID, path, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (m *MinioStorage) ListFiles(bucketID, prefix string) ([]string, error) {
	ctx := context.Background()
	var objects []string

	objectCh := m.Client.ListObjects(ctx, bucketID, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object.Key)
	}
	return objects, nil
}

func (m *MinioStorage) RemoveFile(bucketID string, paths []string) error {
	ctx := context.Background()
	objectsCh := make(chan minio.ObjectInfo, len(paths))

	go func() {
		defer close(objectsCh)
		for _, p := range paths {
			objectsCh <- minio.ObjectInfo{Key: p}
		}
	}()

	opts := minio.RemoveObjectsOptions{}
	for rErr := range m.Client.RemoveObjects(ctx, bucketID, objectsCh, opts) {
		if rErr.Err != nil {
			log.Println("Failed to remove object:", rErr)
		}
	}
	return nil
}

type UploadClient struct {
	storage StorageInterface
}

func NewUploadClient(storage StorageInterface) *UploadClient {
	return &UploadClient{storage: storage}
}

func (u *UploadClient) UploadBuild(ctx context.Context, bucketID, slug string) error {
	buildDir := "./dist"
	files, err := collectFiles(buildDir)
	if err != nil {
		return err
	}

	if err := u.uploadFilesConcurrently(ctx, buildDir, files, bucketID, slug); err != nil {
		_ = u.DeleteDir(slug, bucketID)
		return err
	}
	return nil
}

func (u *UploadClient) UploadFile(baseDir, filePath, bucketID, slug string) error {
	reader, objectKey, err := prepareFileUpload(baseDir, filePath)
	if err != nil {
		return err
	}
	defer reader.Reset(nil)

	return u.storage.UploadFile(bucketID, filepath.ToSlash(filepath.Join(slug, objectKey)), reader)
}

func (u *UploadClient) DeleteDir(dir, bucketID string) error {
	files, err := u.storage.ListFiles(bucketID, dir+"/")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	return u.storage.RemoveFile(bucketID, files)
}

func collectFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (u *UploadClient) uploadFilesConcurrently(ctx context.Context, baseDir string, files []string, bucketID string, slug string) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, fp := range files {
		fp := fp
		g.Go(func() error {
			return u.UploadFile(baseDir, fp, bucketID, slug)
		})
	}
	return g.Wait()
}

func prepareFileUpload(baseDir, filePath string) (*bufio.Reader, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}

	absFile, err := filepath.Abs(filePath)
	if err != nil {
		file.Close()
		return nil, "", err
	}

	rel, err := filepath.Rel(baseDir, absFile)
	if err != nil {
		file.Close()
		return nil, "", err
	}

	reader := bufio.NewReaderSize(file, 1024*1024)
	return reader, rel, nil
}
