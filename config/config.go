package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/juju/errors"
	"strings"
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
	TransferService *TransferService `yaml:"transfer_service"`
}

type TransferService struct {
	Password string `yaml:"password"`
	DropboxDirMapping map[string]string `yaml:"dropbox_dir_mapping"`
}

func (ts TransferService) MapDropboxPath(path string) string {

	for from, to := range ts.DropboxDirMapping {
		if strings.HasPrefix(path, from) {
			return to + path[len(from):]
		}
	}

	return path
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
