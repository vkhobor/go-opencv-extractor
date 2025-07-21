package config

type ServerConfig struct {
	Port        int    `koanf:"port"`
	BlobStorage string `koanf:"blobstorage"`
	BaseUrl     string `koanf:"baseurl"`
	LogFolder   string `koanf:"logfolder"`
	Db          string `koanf:"db"`
}

func (c ServerConfig) GetDirectoryConfig() (DirectoryConfig, error) {
	return newDirectoryConfig(c.BlobStorage)
}
