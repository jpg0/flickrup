package tags

import (
	"github.com/jpg0/flickrup/processing"
	"regexp"
)


type Rewriter struct {
	re *regexp.Regexp
}

func NewRewriter() *Rewriter {
	return &Rewriter{
		re: regexp.MustCompile("^([^:]+):([^:]+)::"),
	}
}

func (rw *Rewriter) MaybeRewrite(ctx *processing.ProcessingContext) processing.ProcessingResult {

	for _, keyword := range ctx.File.Keywords().All().Slice() {
		updated := rw.re.ReplaceAllString(keyword, "$1:$2=")

		if keyword != updated {
			ctx.File.Keywords().Replace(keyword, updated)
			ctx.FileUpdated = true
		}
	}

	return processing.NewSuccessResult()
}