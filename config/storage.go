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
	err := pathutils.EnsurePath(baseDir, true)

	if err != nil {
		return DirectoryConfig{}, err
	}

	return DirectoryConfig{
		BaseDir: baseDir,
	}, nil
}

func (c DirectoryConfig) GetImagesDir() string {
	specific := path.Join(c.BaseDir, imageFolder)
	pathutils.MustEnsurePath(specific, true)
	return specific
}
func (c DirectoryConfig) GetVideosDir() string {
	specific := path.Join(c.BaseDir, videosFolder)
	pathutils.MustEnsurePath(specific, true)
	return specific
}
func (c DirectoryConfig) GetReferencesDir() string {
	specific := path.Join(c.BaseDir, referencesFolder)
	pathutils.MustEnsurePath(specific, true)
	return specific
}
