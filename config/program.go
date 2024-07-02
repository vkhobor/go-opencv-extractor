package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ProgramConfig struct {
	Port        int    `mapstructure:"port"`
	BlobStorage string `mapstructure:"blob_storage"`
	BaseUrl     string `mapstructure:"base_url"`
	Db          string `mapstructure:"db"`
}

func addDefaults(viperConf *viper.Viper) {
	viperConf.SetDefault("port", 8080)
	viperConf.SetDefault("base_url", "http://localhost:8080")
	viperConf.SetDefault("blob_storage", "~/.go_extractor/data")
	viperConf.SetDefault("db", "~/.go_extractor/db.sqlite3")
}

func (c ProgramConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("invalid port")
	}

	if c.BlobStorage == "" {
		return fmt.Errorf("invalid storage path")
	}

	if c.Db == "" {
		return fmt.Errorf("invalid db file")
	}

	return nil
}

var ConfigPaths []string = []string{
	"/etc/go_extractor/",
	"$HOME/.go_extractor/",
	".",
}

func MustNewDefaultViperConfig() *viper.Viper {
	v := viper.New()

	addDefaults(v)
	viper.SetEnvPrefix("GO_EXTRACT")
	viper.AutomaticEnv()
	viper.SetConfigName("config")

	for _, path := range ConfigPaths {
		viper.AddConfigPath(path)
	}

	_, err := NewDefaultProgramConfig(*v)

	if err != nil {
		panic(err)
	}

	return v
}

func NewDefaultProgramConfig(viperConf viper.Viper) (config ProgramConfig, error error) {
	viperConf.Unmarshal(&config)

	if err := config.Validate(); err != nil {
		error = err
	}

	return
}
