package listen

import "github.com/jpg0/flickrup/processing"

func NotifyStage(notifier chan<- struct{}) func(ctx *processing.ProcessingContext, next processing.Processor) processing.ProcessingResult {
	return func(ctx *processing.ProcessingContext, next processing.Processor) processing.ProcessingResult {

		defer func(){notifier <-struct {}{}}()

		return next(ctx)
	}
}