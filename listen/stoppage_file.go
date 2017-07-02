package listen

import (
	"os"
	"io/ioutil"
	"strings"
	"fmt"
	"github.com/Sirupsen/logrus"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"golang.org/x/image/font"
	"image/draw"
	"github.com/golang/freetype/truetype"
	_ "image/jpeg"
	"github.com/golang/freetype"
	"github.com/juju/errors"
	"image/jpeg"
)

const STATUS_FILE_PREFIX = "_flickrup_stoppage_"
var font_files = [...]string{
	"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", //Linux
	"/Library/Fonts/Arial.ttf", //OSX
}

type UploadStatus struct {
	dir string
	cm *ChangeManger
}

func NewUploadStatus(dir string, cm *ChangeManger) *UploadStatus {
	return &UploadStatus{dir:dir, cm:cm}
}

func (us *UploadStatus) IsStatusFile(path string) bool {
	return strings.HasPrefix(path, fmt.Sprintf("%v%v%v", us.dir, string(os.PathSeparator), STATUS_FILE_PREFIX))
}

func (us *UploadStatus) UpdateStatus(filename string) (err error) {
	err = us.ClearStatus()

	if err != nil {
		return
	}

	return us.WriteStatus(filename)
}

func (us *UploadStatus) WriteStatus(filename string) (err error) {
	reader, err := os.Open(fmt.Sprintf("%v%v%v", us.dir, string(os.PathSeparator), filename))
	if err != nil {
		logrus.Error(err)
	}
	defer reader.Close()

	logrus.Debugf("Loading and decoding image")

	img, _, err := image.Decode(reader)

	if err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("Image decoded")

	var drawImg *image.RGBA
	if (img != nil && img.Bounds() != image.ZR) {
		drawImg = image.NewRGBA(img.Bounds())
		draw.Draw(drawImg, img.Bounds(), img, image.Point{0,0}, draw.Src)
	} else {
		drawImg = image.NewRGBA(image.Rect(0, 0, 300, 100))
	}

	col := color.RGBA{255, 0, 255, 255}

	var x = 10 //img.Bounds().Size().X / 2
	var y = drawImg.Bounds().Size().Y / 2

	point := fixed.P(x, y)

	face, err := GetFace(drawImg)

	if err != nil {
		logrus.Error(err)
	}

	d := &font.Drawer{
		Dst:  drawImg,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(filename)

	fout, err := os.Create(fmt.Sprintf("%v%v%v", us.dir, string(os.PathSeparator), STATUS_FILE_PREFIX + filename))
	if err != nil {
		logrus.Error(err)
		return
	}
	defer fout.Close()

	logrus.Debugf("Encoding new image...")
	if err := jpeg.Encode(fout, drawImg, &jpeg.Options{Quality:10}); err != nil {
		logrus.Error(err)
	}

	us.cm.Expect(fout.Name())

	logrus.Debugf("Created file %v", fout.Name())
	return
}

func GetFace(rgba *image.RGBA) (font.Face, error) {

	var fontfile string

	for _, tryfontfile := range font_files {
		_, err := os.Stat(tryfontfile);
		if(os.IsNotExist(err)) {
			continue
		} else {
			fontfile = tryfontfile
			break
		}
	}

	if fontfile == "" {
		return nil, errors.Errorf("Cannot load any font files")
	}

	logrus.Debugf("Loading fontfile %q\n", fontfile)
	b, err := ioutil.ReadFile(fontfile)
	if err != nil {
		logrus.Warn(err)
		return nil, errors.Annotate(err, "Failed to read font file")
	}
	f, err := truetype.Parse(b)
	if err != nil {
		logrus.Warn(err)
		return nil, errors.Annotate(err, "Failed to parse font file")
	}

	// Freetype context
	c := freetype.NewContext()
	c.SetDPI(400)
	c.SetFont(f)
	c.SetFontSize(32)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingNone)

	opts := truetype.Options{}
	opts.Size = float64(rgba.Bounds().Dx() / 10)
	face := truetype.NewFace(f, &opts)

	return face, nil
}

func (us *UploadStatus) ClearStatus() (err error) {
	files, err := ioutil.ReadDir(us.dir)

	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), STATUS_FILE_PREFIX) {
			err = os.Remove(fmt.Sprintf("%v%v%v", us.dir, string(os.PathSeparator), file.Name()))
			if err != nil {
				logrus.Warnf("Failed to remove file %v", file.Name())
			} else {
				logrus.Debugf("Removed file %v", file.Name())
			}
		}
	}

	return
}