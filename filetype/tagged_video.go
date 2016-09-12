package filetype

import (
	"path/filepath"
	"time"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
	"strings"
)

type TaggedVideo struct {
	filepath string
	picasaIni *PicasaIni
}

type TaggedVideoKeywords struct {
	v TaggedVideo
}

func (ti TaggedVideo) Name() string {
	return filepath.Base(ti.filepath)
}

func (ti TaggedVideo) Filepath() string {
	return ti.filepath
}

func (ti TaggedVideo) StringTag(name string) string {
	k, err := ti.picasaIni.cached.GetKey(name)

	if err != nil {
		return ""
	}

	return k.String()
}

func (ti TaggedVideo) Keywords() processing.Keywords {
	return &TaggedVideoKeywords{v: ti}
}

func (ti TaggedVideo) DateTaken() time.Time {
	panic("Not implemented")
}

func (tik TaggedVideoKeywords) All() *processing.TagSet {
	k, err := tik.v.picasaIni.cached.GetKey("keywords")

	if err != nil {
		return processing.NewEmptyTagSet()
	}

	return processing.NewTagSet(k.Strings(","))
}

func (tik TaggedVideoKeywords) Replace(old string, new string) error {
	k, err := tik.v.picasaIni.cached.GetKey("keywords")

	if err != nil {
		return errors.Trace(err)
	}

	all := k.Strings(",")

	for i, str := range all {
		if str == old {
			all[i] = new
		}
	}

	k.SetValue(strings.Join(all, ","))

	return tik.v.picasaIni.ini.SaveTo(tik.v.filepath)
}

func (tik TaggedVideo) ReplaceStringTag(old string, new string) error {
	k, err := tik.picasaIni.cached.GetKey(old)

	if err != nil {
		return errors.Trace(err)
	}

	k.SetValue(new)

	return tik.picasaIni.ini.SaveTo(tik.filepath)}

func NewTaggedVideo(filepath string) (processing.TaggedFile, error) {

	picasa, err := LoadPicasa(filepath)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return &TaggedVideo{
		picasaIni: picasa,
		filepath: filepath,
	}, nil
}

func TaggedVideoFactory() *processing.TaggedFileFactory {
	return &processing.TaggedFileFactory{
		Constructors: map[string]func (filepath string) (processing.TaggedFile, error) {
			"mpg": NewTaggedVideo,
			"mpeg": NewTaggedVideo,
		},
	}
}