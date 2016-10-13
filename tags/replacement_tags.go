package tags

import (
	"github.com/jpg0/flickrup/processing"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

func MaybeReplace(ctx *processing.ProcessingContext) processing.ProcessingResult {
	replacements := ctx.Config.TagReplacements
	allKeywords := ctx.File.Keywords().All()

	if replacements != nil {
		for tagName, tagValueReplacements := range replacements {
			//existing := ctx.File.StringTag(tagName)

			for tagPresent, value := range tagValueReplacements {

				if !strings.HasPrefix(tagPresent, "$") {
					panic("Static tag replacements not supported")
				}

				if allKeywords.Contains(tagPresent[1:]) {
					log.Debugf("Setting tag %v to %v", tagName, value)

					err := ctx.File.ReplaceStringTag(tagName, value)

					if err != nil {
						return processing.NewErrorResult(errors.Annotate(err, "Replacing string tag"))
					}

					break
				}

			}
		}
	}

	return processing.NewSuccessResult()
}