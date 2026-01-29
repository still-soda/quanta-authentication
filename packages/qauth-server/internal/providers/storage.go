package providers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"qauth-server/internal/config"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

type IStorage interface {
	Upload(ctx context.Context, key string, data io.Reader) error
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetURL(ctx context.Context, key string) (string, error)
}

type LocalStorage struct {
	bucket     *blob.Bucket
	localDir   string
	bucketName string
}

// NewLocalStorage 创建存储服务
func NewLocalStorage(cfg *config.Config) *LocalStorage {
	ctx := context.Background()
	var bucket *blob.Bucket

	localDir := cfg.Storage.LocalDir

	if err := os.MkdirAll(localDir, 0755); err != nil {
		panic("failed to create storage directory: " + err.Error())
	}

	absPath, err := filepath.Abs(localDir)
	if err != nil {
		panic("failed to get absolute path: " + err.Error())
	}

	bucket, err = blob.OpenBucket(ctx, fmt.Sprintf("file://%s", absPath))
	if err != nil {
		panic("failed to open local bucket: " + err.Error())
	}

	return &LocalStorage{
		bucket:     bucket,
		localDir:   absPath,
		bucketName: "local",
	}
}

// Upload 上传文件
func (s *LocalStorage) Upload(ctx context.Context, key string, data io.Reader) error {
	w, err := s.bucket.NewWriter(ctx, key, nil)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}

	if _, err := io.Copy(w, data); err != nil {
		w.Close()
		return fmt.Errorf("failed to write data: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

// Download 下载文件
func (s *LocalStorage) Download(ctx context.Context, key string) ([]byte, error) {
	r, err := s.bucket.NewReader(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return data, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	if err := s.bucket.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	return s.bucket.Exists(ctx, key)
}

// GetURL 获取文件访问 URL
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	return filepath.Join("/uploads", key), nil
}

// Close 关闭存储桶
func (s *LocalStorage) Close() error {
	if s.bucket != nil {
		return s.bucket.Close()
	}
	return nil
}
