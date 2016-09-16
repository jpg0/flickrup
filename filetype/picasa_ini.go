package filetype

import (
	"gopkg.in/ini.v1"
	"path/filepath"
	"github.com/juju/errors"
)

type PicasaIni struct {
	ini *ini.File
	cached *ini.Section
}

func LoadPicasa(path string) (*PicasaIni, error) {
	cfg, err := ini.Load(filepath.Dir(path) + "/.picasa.ini")

	if err != nil {
		return nil, errors.Annotate(err, "Loading Picasa config")
	}

	section := cfg.Section(filepath.Base(path))

	return &PicasaIni{
		ini: cfg,
		cached: section,
	}, nil
}


