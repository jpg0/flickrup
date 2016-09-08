package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/juju/errors"
)

type Config struct {
	WatchDir string `yaml:"watch_dir"`
	APIKey string `yaml:"api_key"`
	SharedSecret string `yaml:"shared_secret"`
	ArchiveDir string `yaml:"archive_dir"`
	TagsetPrefix string `yaml:"tagsetprefix"`
	VisibilityPrefix string `yaml:"visibilityprefix"`
	TagReplacements map[string]map[string]string `yaml:"tag_replacements"`
	BlockedTags map[string]string `yaml:"blocked_tags"`
	ConvertFiles map[string]string `yaml:"convert_files"`
	TransferServicePassword string `yaml:"transfer_service_password"`
}

func Load(filepath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, errors.Trace(err)
	}

	rv := new(Config)

	err = yaml.Unmarshal(bytes, rv)

	return rv, err
}
