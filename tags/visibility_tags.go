package tags

import (
	"github.com/jpg0/flickrup/processing"
	log "github.com/Sirupsen/logrus"
)

func ExtractVisibility(ctx *processing.ProcessingContext) processing.ProcessingResult {
	prefix := ctx.Config.VisibilityPrefix

	if prefix != "" {
		visibilities := processing.ValuesByPrefix(ctx.File.Keywords(), prefix)

		if len(visibilities) == 0 {
			log.Debugf("No visibility specified for %v, using default", ctx.File.Name())
		} else {
			if len(visibilities) > 1 {
				log.Warnf("Multiple visibilities specified for %v, using the first only", ctx.File.Name())
			}

			ctx.Visibilty = visibilities[0]
			log.Debugf("Updated visibility to %v for %v", ctx.Visibilty, ctx.File.Name())
		}
	}

	return processing.NewSuccessResult()
}