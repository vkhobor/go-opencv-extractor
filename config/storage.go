package config

import (
	"path"

	pathutils "github.com/vkhobor/go-opencv/path"
)

const imageFolder = "images"
const videosFolder = "videos"
const referencesFolder = "references"

type DirectoryConfig struct {
	BaseDir string
}

func NewDirectoryConfig(baseDir string) (DirectoryConfig, error) {
	path, err := pathutils.EnsurePath(baseDir)

	if err != nil {
		return DirectoryConfig{}, err
	}

	return DirectoryConfig{
		BaseDir: path,
	}, nil
}

func (c DirectoryConfig) GetImagesDir() string {
	return path.Join(c.BaseDir, imageFolder)
}
func (c DirectoryConfig) GetVideosDir() string {
	return path.Join(c.BaseDir, videosFolder)
}
func (c DirectoryConfig) GetReferencesDir() string {
	return path.Join(c.BaseDir, referencesFolder)
}
