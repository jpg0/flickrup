package filetype

import (
	"path/filepath"
	"time"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
	"strings"
	"os"
	"github.com/jpg0/goexiftool"
	"github.com/Sirupsen/logrus"
)

type TaggedVideo struct {
	filepath string
	picasaIni *PicasaIni
	dateTaken time.Time
	img goexiftool.Image
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
	return ti.dateTaken
}

func (tik TaggedVideoKeywords) All() (rv *processing.TagSet) {
	if tik.v.picasaIni != nil {
		k, err := tik.v.picasaIni.cached.GetKey("keywords")

		if err != nil {
			logrus.Debugf("No picasa tags loaded for %s [%s]", tik.v.filepath, err)
			rv = processing.NewEmptyTagSet()
		} else {
			rv = processing.NewTagSet(k.Strings(","))
		}
	}

	if tik.v.img != nil {
		//and EXIF tags
		tags, err := tik.v.img.StringSlice("Keywords")

		if err != nil {
			logrus.Debugf("No exif tags [Keywords] loaded for %s [%s]", tik.v.filepath, err)
		} else {
			rv.AddAll(processing.NewTagSet(tags))
		}

		tags, err = tik.v.img.StringSlice("Subject")

		if err != nil {
			logrus.Debugf("No exif tags [Subject] loaded for %s [%s]", tik.v.filepath, err)
		} else {
			rv.AddAll(processing.NewTagSet(tags))
		}
	}

	return
}

func (tik TaggedVideoKeywords) Replace(old string, new string) (err error) {

	//first picasa ini
	if tik.v.picasaIni != nil {
		k, err := tik.v.picasaIni.cached.GetKey("keywords")

		if err != nil {
			logrus.Errorf("Failed to load (& update) picasa tags for %s: %s", tik.v.filepath, err)
		} else {
			all := k.Strings(",")

			if len(all) > 0 {
				for i, str := range all {
					if str == old {
						all[i] = new
					}
				}

				k.SetValue(strings.Join(all, ","))

				err = tik.v.picasaIni.ini.SaveTo(tik.v.picasaIni.filepath)

				if err != nil {
					logrus.Errorf("Failed to update picasa tags for %s: %s", tik.v.filepath, err)
				}
			}
		}
	}

	//and exif
	if tik.v.img != nil {
		//and EXIF tags
		all, err := tik.v.img.StringSlice("Keywords")

		if err != nil {
			logrus.Debugf("No exif tags [Keywords] loaded for %s [%s]", tik.v.filepath, err)
		} else {
			if len(all) > 0 {
				for i, str := range all {
					if str == old {
						all[i] = new
					}
				}

				err = tik.v.img.AddTag("Keywords", strings.Join(all, ","))

				if err != nil {
					logrus.Errorf("Failed to update exif tags for %s: %s", tik.v.filepath, err)
				}
			}
		}
	}

	return
}

func (tik TaggedVideo) ReplaceStringTag(old string, new string) (err error) {

	//first picasa ini
	if tik.picasaIni != nil {
		k, err := tik.picasaIni.cached.GetKey("keywords")

		if err != nil {
			logrus.Errorf("Failed to load (& update) picasa tags for %s: %s", tik.filepath, err)
		} else {
			k.SetValue(new);
			err = tik.picasaIni.ini.SaveTo(tik.picasaIni.filepath)

			if err != nil {
				logrus.Errorf("Failed to update picasa tags for %s: %s", tik.filepath, err)
			}
		}
	}

	//and exif
	if tik.img != nil {
		//and EXIF tags

		err := tik.img.AddTag("Keywords", new)

		if err != nil {
			logrus.Errorf("Failed to update exif tags for %s: %s", tik.filepath, err)
		}
	}

	return
}

func NewTaggedVideo(filepath string) (processing.TaggedFile, error) {

	picasa, picasaErr := LoadPicasa(filepath)

	if picasaErr != nil {
		logrus.Warnf("Failed to read picasa config: %s", picasaErr)
	}

	img, imgErr := goexiftool.NewImage(filepath)

	if imgErr != nil {
		logrus.Warnf("Failed to read video EXIF for %s: %s", filepath, imgErr)
	}

	if imgErr != nil && picasaErr != nil {
		return nil, errors.Annotate(imgErr, "Reading video EXIF")
	}

	dateTaken, err := os.Stat(filepath)

	if err != nil {
		return nil, errors.Annotate(err, "Reading file info for date")
	}

	return &TaggedVideo{
		picasaIni: picasa,
		img: img,
		filepath: filepath,
		dateTaken: dateTaken.ModTime(),
	}, nil
}

func TaggedVideoFactory() *processing.TaggedFileFactory {
	return &processing.TaggedFileFactory{
		Constructors: map[string]func (filepath string) (processing.TaggedFile, error) {
			"mov": NewTaggedVideo,
			"mpeg": NewTaggedVideo,
			"mp4": NewTaggedVideo,
		},
	}
}
