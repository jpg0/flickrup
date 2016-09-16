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
	rv, ok := ti.img.Tags()[name].(string)

	if ok {
		return rv
	} else {
		return ""
	}
}

func (ti TaggedImage) Keywords() processing.Keywords {
	return TaggedImageKeywords{img: ti}
}

func (ti TaggedImage) DateTaken() time.Time {
	rv := ti.img.Tags()["DateTimeOriginal"]

	if rv == "" {
		rv = ti.img.Tags()["ModifyDate"]
	}

	if rv == "" {
		rv = ti.img.Tags()["FileModifyDate"]
	}

	timeString, ok := rv.(string)

	if ok {
		t, e := time.Parse(TIME_FORMAT, timeString)

		if e != nil {
			panic("Failed to parse time: " + rv.(string))
		}

		return t
	}

	timeUint64, ok := rv.(int64)

	if ok {
		if timeUint64 > 999999999 { //in millis
			return time.Unix(timeUint64 / 1000, timeUint64 % 1000)
		} else {
			return time.Unix(timeUint64, 0)
		}
	}

	return time.Time{}
}

func (tik TaggedImageKeywords) All() *processing.TagSet {

	kw, err := tik.img.img.StringSlice("Keywords")

	if err != nil {
		panic("Failure reading keywords: " + err.Error())
	}

	sub, err := tik.img.img.StringSlice("Keywords")

	if err != nil {
		panic("Failure reading keywords: " + err.Error())
	}

	ts := processing.NewTagSet(kw)
	ts.AddAll(processing.NewTagSet(sub))

	return ts
}

func (tik TaggedImageKeywords) Replace(old string, new string) error {

	err := tik.writeKeywordReplacement(old, new)

	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (tik TaggedImageKeywords) writeKeywordReplacement(old string, new string) error {
	err := tik.img.img.RemoveTagValue("Keywords", old)
	if err != nil { return errors.Trace(err) }

	err = tik.img.img.AddTagValue("Keywords", new)
	if err != nil { return errors.Trace(err) }

	err = tik.img.img.RemoveTagValue("Subject", old)
	if err != nil { return errors.Trace(err) }

	err = tik.img.img.AddTagValue("Subject", new)
	if err != nil { return errors.Trace(err) }

	return nil
}

func (tik TaggedImage) ReplaceStringTag(name string, newValue string) error {
	err := tik.img.RemoveTag(name)

	if err != nil {
		return errors.Trace(err)
	}

	return tik.img.AddTag(name, newValue)
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