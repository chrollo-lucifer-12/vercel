package upload

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/chrollo-lucifer-12/shared/utils"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/sync/errgroup"
)

type MockStorage struct {
	UploadFileFn func(bucketID, path string, reader io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error)
	ListFilesFn  func(bucketID, path string, opts storage_go.FileSearchOptions) ([]storage_go.FileObject, error)
	RemoveFileFn func(bucketID string, paths []string) ([]storage_go.FileUploadResponse, error)
}

func (m *MockStorage) UploadFile(bucketID, path string, reader io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error) {
	if m.UploadFileFn != nil {
		return m.UploadFileFn(bucketID, path, reader, fileOptions...)
	}
	return storage_go.FileUploadResponse{}, nil
}

func (m *MockStorage) ListFiles(bucketID, path string, opts storage_go.FileSearchOptions) ([]storage_go.FileObject, error) {
	if m.ListFilesFn != nil {
		return m.ListFilesFn(bucketID, path, opts)
	}
	return nil, nil
}

func (m *MockStorage) RemoveFile(bucketID string, paths []string) ([]storage_go.FileUploadResponse, error) {
	if m.RemoveFileFn != nil {
		return m.RemoveFileFn(bucketID, paths)
	}
	return nil, nil
}

type StorageInterface interface {
	UploadFile(bucketID, path string, reader io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error)
	ListFiles(bucketID, path string, opts storage_go.FileSearchOptions) ([]storage_go.FileObject, error)
	RemoveFile(bucketID string, paths []string) ([]storage_go.FileUploadResponse, error)
}

type UploadClient struct {
	storage StorageInterface
}

func NewUploadClient(apiURL, apiKey string) (*UploadClient, error) {
	client, err := supabase.NewClient(apiURL, apiKey, nil)
	if err != nil {
		return nil, err
	}
	return &UploadClient{storage: client.Storage}, nil
}

func NewUploadClientWithStorage(storage StorageInterface) *UploadClient {
	return &UploadClient{storage: storage}
}

func (u *UploadClient) UploadBuild(ctx context.Context, bucketID, slug string) error {
	buildDir, err := u.getBuildDirectory()
	if err != nil {
		return err
	}

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

	_, err = u.storage.UploadFile(bucketID, filepath.ToSlash(filepath.Join(slug, objectKey)), reader)
	return err
}

func (u *UploadClient) DeleteDir(dir, bucketID string) error {
	files, err := u.storage.ListFiles(bucketID, dir+"/", storage_go.FileSearchOptions{
		Limit: 1000,
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	paths := make([]string, 0, len(files))
	for _, f := range files {
		paths = append(paths, f.Name)
	}

	_, err = u.storage.RemoveFile(bucketID, paths)
	return err
}

func (u *UploadClient) getBuildDirectory() (string, error) {
	buildPath := utils.GetPath([]string{"home", "app", "output", "dist"})
	return filepath.Abs(buildPath)
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

func (u *UploadClient) uploadFilesConcurrently(
	ctx context.Context,
	baseDir string,
	files []string,
	bucketID string,
	slug string,
) error {

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
