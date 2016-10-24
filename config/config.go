package config

import (
	//"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/juju/errors"
	"strings"
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

type Config struct {
	WatchDir string `json:"watch_dir"`
	APIKey string `json:"api_key"`
	SharedSecret string `json:"shared_secret"`
	ArchiveDir string `json:"archive_dir"`
	TagsetPrefix string `json:"tagsetprefix"`
	VisibilityPrefix string `json:"visibilityprefix"`
	TagReplacements map[string]map[string]string `json:"tag_replacements"`
	BlockedTags map[string]string `json:"blocked_tags"`
	ConvertFiles map[string][]string `json:"convert_files"`
	TransferService *TransferService `json:"transfer_service"`
}

type TransferService struct {
	Password string `json:"password"`
	DropboxDirMapping map[string]string `json:"dropbox_dir_mapping"`
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

	err = json.Unmarshal(bytes, rv)

	if err != nil {
		logrus.Debugf("Loaded config %s", filepath)
	}

	return rv, err
}
