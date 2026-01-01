package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"nexushub-personal/internal/config"
	"nexushub-personal/internal/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type CloudStorageProvider interface {
	// 上传文件
	Upload(ctx context.Context, objectName string, file []byte, contentType string) error
	// 下载文件
	Download(ctx context.Context, objectName string) ([]byte, error)
	// 获取文件URL
	GetFileURL(ctx context.Context, objectName string) (string, error)
	// 删除文件
	Delete(ctx context.Context, objectName string) error
}

// MinIOProvider MinIO云存储实现
type MinIOProvider struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOProvider 创建MinIO云存储实例
func NewMinIOProvider() (*MinIOProvider, error) {
	cfg := config.AppConfig.GetCloudStorageConfig()

	if cfg.Provider != string(config.ProviderMinIO) {
		return nil, fmt.Errorf("provider mismatch: expected minio, got %s", cfg.Provider)
	}

	// 创建MinIO客户端
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Endpoint != "", // 如果配置了endpoint，则使用HTTPS
	})

	if err != nil {
		logger.Error("Failed to create MinIO client: %v", err)
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// 检查存储桶是否存在，不存在则创建
	err = client.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		// 如果是存储桶已存在的错误，则忽略
		exists, errBucketExists := client.BucketExists(context.Background(), cfg.Bucket)
		if errBucketExists != nil || !exists {
			logger.Error("Failed to create bucket: %v", err)
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinIOProvider{
		client:     client,
		bucketName: cfg.Bucket,
	}, nil
}

// Upload 上传文件到MinIO
func (m *MinIOProvider) Upload(ctx context.Context, objectName string, file []byte, contentType string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, bytes.NewReader(file), int64(len(file)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		logger.Error("Failed to upload file to MinIO: %v, object: %s", err, objectName)
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return nil
}

// Download 从MinIO下载文件
func (m *MinIOProvider) Download(ctx context.Context, objectName string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("Failed to download file from MinIO: %v, object: %s", err, objectName)
		return nil, fmt.Errorf("failed to download file from MinIO: %w", err)
	}
	defer object.Close()

	// 获取文件大小
	stat, err := object.Stat()
	if err != nil {
		logger.Error("Failed to get file info from MinIO: %v, object: %s", err, objectName)
		return nil, fmt.Errorf("failed to get file info from MinIO: %w", err)
	}

	// 读取文件内容
	fileBytes := make([]byte, stat.Size)
	_, err = object.Read(fileBytes)
	if err != nil && err != io.EOF {
		logger.Error("Failed to read file from MinIO: %v, object: %s", err, objectName)
		return nil, fmt.Errorf("failed to read file from MinIO: %w", err)
	}

	return fileBytes, nil
}

// GetFileURL 获取文件访问URL
func (m *MinIOProvider) GetFileURL(ctx context.Context, objectName string) (string, error) {
	cfg := config.AppConfig.GetCloudStorageConfig()

	// 如果配置了自定义域名或endpoint，使用该地址
	if cfg.Endpoint != "" {
		return cfg.Endpoint + "/" + m.bucketName + "/" + objectName, nil
	}

	// 否则使用默认的MinIO URL格式
	// 注意：这里假设MinIO服务器地址是通过配置文件获取的
	// 在实际项目中，您可能需要从配置中获取服务器地址
	// 为了演示，我们使用一个简化的URL生成方式
	return "", fmt.Errorf("endpoint not configured for MinIO")
}

// Delete 从MinIO删除文件
func (m *MinIOProvider) Delete(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})

	if err != nil {
		logger.Error("Failed to delete file from MinIO: %v, object: %s", err, objectName)
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}
