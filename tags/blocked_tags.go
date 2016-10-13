package tags

import ("github.com/jpg0/flickrup/processing"
	"errors"
	"fmt"
)

func MaybeBlock(ctx *processing.ProcessingContext) processing.ProcessingResult {
	blockers := ctx.Config.BlockedTags

	if blockers != nil {
		for tag, value := range blockers {
			if ctx.File.StringTag(tag) == value {
				return processing.NewErrorResult(errors.New(fmt.Sprintf("Blocking upload on %v as %v=%v", ctx.File.Name(), tag, value)))
			}
		}
	}

	return processing.NewSuccessResult()
}