package filetype

import (
	"github.com/jpg0/goexiftool"
	"path/filepath"
	"time"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
)

const TIME_FORMAT = "2006:01:02 15:04:05"

type TaggedImage struct {
	filepath string
	img goexiftool.Image
}

type TaggedImageKeywords struct {
	img TaggedImage
}

func (ti TaggedImage) Name() string {
	return filepath.Base(ti.filepath)
}

func (ti TaggedImage) Filepath() string {
	return ti.filepath
}

func (ti TaggedImage) StringTag(name string) string {
	return ti.img.Tags()[name].(string)
}

func (ti TaggedImage) Keywords() processing.Keywords {
	return TaggedImageKeywords{img: ti}
}

func (ti TaggedImage) RealDateTaken() time.Time {
	rv := ti.img.Tags()["DateTimeOriginal"]

	if rv == "" {
		rv = ti.img.Tags()["ModifyDate"]
	}

	if rv == "" {
		rv = ti.img.Tags()["FileModifyDate"]
	}

	t, e := time.Parse(TIME_FORMAT, rv.(string))

	if e != nil {
		panic("Failed to parse time: " + rv.(string))
	}

	return t
}

func (tik TaggedImageKeywords) All() []string {
	kw := tik.img.img.Tags()["Keywords"]

	if s, ok := kw.(string); ok {
		return []string{s}
	} else if ss, ok := kw.([]interface{}); ok { //need to assert via []interface{}
		rv := make([]string, len(ss))
		for i, s := range ss { rv[i] = s.(string) }
		return rv
	} else { //assume unset
		return []string{}
	}
}

func (tik TaggedImageKeywords) Replace(old string, new string) error {
	panic("no implemented")
}

func (tik TaggedImage) ReplaceStringTag(old string, new string) error {
	panic("no implemented")
}

func NewTaggedImage(filepath string) (processing.TaggedFile, error) {
	img, err := goexiftool.NewImage(filepath)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return &TaggedImage{
		filepath: filepath,
		img: img,
	}, nil
}

func TaggedImageFactory() *processing.TaggedFileFactory {
	return &processing.TaggedFileFactory{
		Constructors: map[string]func (filepath string) (processing.TaggedFile, error) {
			"jpg": NewTaggedImage,
			"jpeg": NewTaggedImage,
		},
	}
}