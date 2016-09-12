package filetype

import (
	"github.com/jpg0/flickrup/processing"
	"io/ioutil"
	"fmt"
	"strings"
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
	return ioutil.WriteFile(archived + ".meta", []byte(data), 0644)
}