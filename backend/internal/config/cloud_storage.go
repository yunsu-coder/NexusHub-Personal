package config

type CloudStorageConfig struct {
	Provider  string
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Region    string
}

type CloudProvider string

const (
	ProviderMinIO  CloudProvider = "minio"
	ProviderAWS    CloudProvider = "aws"
	ProviderAliyun CloudProvider = "aliyun"
	ProviderTencent CloudProvider = "tencent"
)

func (c *Config) GetCloudStorageConfig() *CloudStorageConfig {
	return &CloudStorageConfig{
		Provider:  getEnv("CLOUD_STORAGE_PROVIDER", ""),
		AccessKey: getEnv("CLOUD_STORAGE_ACCESS_KEY", ""),
		SecretKey: getEnv("CLOUD_STORAGE_SECRET_KEY", ""),
		Bucket:    getEnv("CLOUD_STORAGE_BUCKET", ""),
		Endpoint:  getEnv("CLOUD_STORAGE_ENDPOINT", ""),
		Region:    getEnv("CLOUD_STORAGE_REGION", ""),
	}
}

// GetCloudStorageURL 获取云存储的文件访问URL
func (c *Config) GetCloudStorageURL(filename string) string {
	storage := c.GetCloudStorageConfig()
	
	// 如果未配置云存储，返回空
	if storage.Provider == "" {
		return ""
	}
	
	// 构建存储路径
	storagePath := ""
	if storage.Bucket != "" {
		storagePath = storage.Bucket + "/"
	}
	
	// 构建完整的URL
	if storage.Endpoint != "" {
		return storage.Endpoint + "/" + storagePath + filename
	}
	
	return ""
}

// IsCloudStorageEnabled 检查是否启用了云存储
func (c *Config) IsCloudStorageEnabled() bool {
	storage := c.GetCloudStorageConfig()
	return storage.Provider != ""
}