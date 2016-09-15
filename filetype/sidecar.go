package filetype

import (
	"github.com/jpg0/flickrup/processing"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/Sirupsen/logrus"
)

func SidecarStage() func(ctx *processing.ProcessingContext, next processing.Processor) error {
	return func(ctx *processing.ProcessingContext, next processing.Processor) error {

		err := next(ctx)

		if err == nil {
			asVideo, ok := ctx.File.(*TaggedVideo)

			if ok {
				return writeSidecar(asVideo, ctx.ArchivedAs)
			}
		}

		return err
	}
}

func writeSidecar(video *TaggedVideo, archived string) error {
	data := fmt.Sprintf("keywords: %v", strings.Join(video.Keywords().All().Slice(), ","))
	sidecarPath := archived + ".meta"

	logrus.Debugf("Writing sidecar as %v", sidecarPath)

	return ioutil.WriteFile(sidecarPath, []byte(data), 0644)
}