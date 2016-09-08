package filetype

import (
	"fmt"
	"github.com/jpg0/goexiftool"
	"path/filepath"
	"time"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
)

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

	if rv == nil {
		rv = ti.img.Tags()["ModifyDate"]
	}

	if rv == nil {
		rv = ti.img.Tags()["FileModifyDate"]
	}

	return rv.(time.Time)
}

func (tik TaggedImageKeywords) All() []string {
	kw := tik.img.img.Tags()["Keywords"]

	switch kw.(type) {
	default:
		panic(fmt.Sprintf("unexpected tag type %T", kw))
	case nil:
		return []string{""}
	case string:
		return []string{kw.(string)}
	case []string:
		return kw.([]string)
	}
}

func (tik TaggedImageKeywords) Replace(old string, new string) error {
	panic("no implemented")
}

func (tik TaggedImage) ReplaceStringTag(old string, new string) error {
	panic("no implemented")
}

func NewTaggedImage(filepath string) (*TaggedImage, error) {
	img, err := goexiftool.NewImage(filepath)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return &TaggedImage{
		filepath: filepath,
		img: img,
	}, nil
}