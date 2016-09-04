package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	WatchDir string `yaml:"watch_dir"`
	APIKey string `yaml:"api_key"`
	SharedSecret string
	ArchiveDir string `yaml:"archive_dir"`
	TagsetPrefix string `yaml:"tagsetprefix"`
	VisibilityPrefix string `yaml:"visibilityprefix"`
	TagReplacements string `yaml:"tag_replacements"`
	LensType map[string]string
	FocalLength map[string]string
	FocalLengthIn35mmFormat map[string]string
	BlockedTags map[string]string `yaml:"blocked_tags"`
	ConvertFiles map[string]string `yaml:"convert_files"`
}

func Load(filepath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	rv := new(Config)

	err = yaml.Unmarshal(bytes, rv)

	return rv, err
}
