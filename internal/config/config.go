package config

import "errors"

type Config struct {
	Storage     StorageConfigInterface
	TokenSecret string
	Postgres    PostgresConfig
}

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Port     int
}

type StorageConfigInterface interface {
	LoadConfig() error
	GetStorageType() string
}

type StorageConfig struct {
	StorageType string
	Local       LocalConfig
	AzureBlob   AzureBlobConfig
}

func (sc *StorageConfig) LoadConfig() error {
	switch sc.StorageType {
	case STORAGE_TYPE_LOCAL:
		return sc.Local.LoadConfig()
	case STORAGE_TYPE_AZURE_BLOB:
		return sc.AzureBlob.LoadConfig()
	default:
		return errors.New("invalid storage type")
	}
}

func (sc *StorageConfig) GetStorageType() string {
	return sc.StorageType
}

type LocalConfig struct {
	Directory string
}

func (lc *LocalConfig) LoadConfig() error {
	return nil
}

func (lc *LocalConfig) GetStorageType() string {
	return STORAGE_TYPE_LOCAL
}

type AzureBlobConfig struct {
	AccountName string
	AccountKey  string
	Container   string
}

func (abc *AzureBlobConfig) LoadConfig() error {
	return nil
}

func (abc *AzureBlobConfig) GetStorageType() string {
	return STORAGE_TYPE_AZURE_BLOB
}

const (
	STORAGE_TYPE_LOCAL      = "local"
	STORAGE_TYPE_AZURE_BLOB = "azure_blob"
)

func LoadConfig() *Config {
	return &Config{
		TokenSecret: "secret",
		Postgres: PostgresConfig{
			Port:     5432,
			Host:     "localhost",
			User:     "postgres",
			Password: "password",
			Name:     "riley",
			SSLMode:  "disable",
		},
		Storage: &StorageConfig{
			StorageType: STORAGE_TYPE_AZURE_BLOB,
			AzureBlob: AzureBlobConfig{
				AccountName: "account",
				AccountKey:  "key",
				Container:   "container",
			},
		},
	}
}

func LoadTestConfig() *Config {
	return &Config{
		TokenSecret: "secret",
		Postgres: PostgresConfig{
			Port:     5432,
			Host:     "localhost",
			User:     "postgres",
			Password: "password",
			Name:     "postgres",
			SSLMode:  "disable",
		},
		Storage: &StorageConfig{
			StorageType: STORAGE_TYPE_LOCAL,
			Local: LocalConfig{
				Directory: "/tmp",
			},
		},
	}
}
