package filetype

import (
	"gopkg.in/ini.v1"
	"path/filepath"
)

type PicasaIni struct {
	ini *ini.File
	cached *ini.Section
}

func LoadPicasa(path string) (*PicasaIni, error) {
	cfg, err := ini.Load(path + "/.picasa.ini")

	if err != nil {
		return nil, err
	}

	section, err := cfg.GetSection(filepath.Base(path))

	if err != nil {
		return nil, err
	}

	return &PicasaIni{
		ini: cfg,
		cached: section,
	}, nil
}


