package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"qauth-server/internal/config"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

// Service 存储服务接口
type Service interface {
	Upload(ctx context.Context, key string, data io.Reader) error
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetURL(ctx context.Context, key string) (string, error)
}

// CDKStorage Go CDK 存储实现
type CDKStorage struct {
	bucket     *blob.Bucket
	provider   string
	localDir   string
	bucketName string
}

// NewCDKStorage 创建新的 CDK 存储服务
func NewCDKStorage(cfg *config.Config) (*CDKStorage, error) {
	ctx := context.Background()
	var bucket *blob.Bucket
	var err error

	provider := cfg.Storage.Provider

	switch provider {
	case "local":
		localDir := cfg.Storage.LocalDir

		if err := os.MkdirAll(localDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create storage directory: %w", err)
		}

		absPath, err := filepath.Abs(localDir)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %w", err)
		}

		bucket, err = blob.OpenBucket(ctx, fmt.Sprintf("file://%s", absPath))
		if err != nil {
			return nil, fmt.Errorf("failed to open local bucket: %w", err)
		}

		log.Printf("✅ 本地存储初始化成功: %s\n", absPath)
		return &CDKStorage{
			bucket:   bucket,
			provider: provider,
			localDir: absPath,
		}, nil

	case "s3":
		bucketName := cfg.Storage.S3Bucket
		region := cfg.Storage.S3Region

		if bucketName == "" {
			return nil, fmt.Errorf("S3 bucket name is required")
		}

		bucketURL := fmt.Sprintf("s3://%s?region=%s", bucketName, region)
		bucket, err = blob.OpenBucket(ctx, bucketURL)
		if err != nil {
			return nil, fmt.Errorf("failed to open S3 bucket: %w", err)
		}

		log.Printf("✅ S3 存储初始化成功: %s (region: %s)\n", bucketName, region)
		return &CDKStorage{
			bucket:   bucket,
			provider: provider,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", provider)
	}
}

// Upload 上传文件
func (s *CDKStorage) Upload(ctx context.Context, key string, data io.Reader) error {
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
func (s *CDKStorage) Download(ctx context.Context, key string) ([]byte, error) {
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
func (s *CDKStorage) Delete(ctx context.Context, key string) error {
	if err := s.bucket.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// Exists 检查文件是否存在
func (s *CDKStorage) Exists(ctx context.Context, key string) (bool, error) {
	return s.bucket.Exists(ctx, key)
}

// GetURL 获取文件访问 URL
func (s *CDKStorage) GetURL(ctx context.Context, key string) (string, error) {
	if s.provider == "local" {
		return filepath.Join("/uploads", key), nil
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key), nil
}

// Close 关闭存储桶
func (s *CDKStorage) Close() error {
	if s.bucket != nil {
		return s.bucket.Close()
	}
	return nil
}
