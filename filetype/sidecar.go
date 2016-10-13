package filetype

import (
	"github.com/jpg0/flickrup/processing"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/Sirupsen/logrus"
)

func SidecarStage() func(ctx *processing.ProcessingContext, next processing.Processor) processing.ProcessingResult {
	return func(ctx *processing.ProcessingContext, next processing.Processor) processing.ProcessingResult {

		result := next(ctx)

		if result.ResultType != processing.SuccessResult {
			asVideo, ok := ctx.File.(*TaggedVideo)

			if ok {
				err := writeSidecar(asVideo, ctx.ArchivedAs)

				if err != nil {
					return processing.NewErrorResult(err)
				} else {
					return processing.NewSuccessResult()
				}
			}
		}

		return result
	}
}

func writeSidecar(video *TaggedVideo, archived string) error {
	data := fmt.Sprintf("keywords: %v", strings.Join(video.Keywords().All().Slice(), ","))
	sidecarPath := archived + ".meta"

	logrus.Debugf("Writing sidecar as %v", sidecarPath)

	return ioutil.WriteFile(sidecarPath, []byte(data), 0644)
}