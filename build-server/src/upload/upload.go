package upload

import (
	"bufio"
	"context"
	"os"
	"path/filepath"

	"github.com/chrollo-lucifer-12/build-server/src/utils"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/sync/errgroup"
)

type UploadClient struct {
	client *supabase.Client
}

func NewUploadClient(apiUrl, apiKey string) (*UploadClient, error) {
	client, err := supabase.NewClient(apiUrl, apiKey, nil)
	if err != nil {
		return nil, err
	}
	return &UploadClient{client: client}, nil
}
func (u *UploadClient) UploadFile(baseDir, filename, bucketID, slug string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	absFile, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(baseDir, absFile)
	if err != nil {
		return err
	}

	objectKey := filepath.ToSlash(rel)
	reader := bufio.NewReaderSize(file, 1024*1024)

	_, err = u.client.Storage.UploadFile(bucketID, slug+"/"+objectKey, reader)
	return err
}

func (u *UploadClient) DeleteDir(dir string, bucketId string) error {
	files, err := u.client.Storage.ListFiles(bucketId, dir+"/", storage_go.FileSearchOptions{
		Limit: 1000,
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = f.Name
	}

	_, err = u.client.Storage.RemoveFile(bucketId, paths)

	return err
}

func (u *UploadClient) UploadBuild(ctx context.Context, bucketID, slug string) error {
	buildPath := utils.GetPath([]string{"home", "app", "output", "dist"})

	absBuildDir, err := filepath.Abs(buildPath)
	if err != nil {
		return err
	}

	files := []string{}
	err = filepath.WalkDir(absBuildDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, fp := range files {
		fp := fp

		g.Go(func() error {
			if err := u.UploadFile(absBuildDir, fp, bucketID, slug); err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {

		u.DeleteDir(slug, bucketID)
		return err
	}

	return nil
}
