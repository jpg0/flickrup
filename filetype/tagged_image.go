package filetype

import (
	"fmt"
	"github.com/jpg0/goexiftool"
	"path/filepath"
)

type TaggedImage struct {
	filepath string
	img goexiftool.Image
}

func (ti TaggedImage) Name() string {
	return filepath.Base(ti.filepath)
}

func (ti TaggedImage) Filepath() string {
	return ti.filepath
}

func (ti TaggedImage) Keywords() []string {
	kw := ti.img.Tags()["Keywords"]

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

func NewTaggedImage(filepath string) (*TaggedImage, error) {
	img, err := goexiftool.NewImage(filepath)

	if err != nil {
		return nil, err
	}

	return &TaggedImage{
		filepath: filepath,
		img: img,
	}, nil
}