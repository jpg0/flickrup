package processing

type Stage func(ctx *ProcessingContext, next Processor) ProcessingResult

func SuccessStage(_ *ProcessingContext) ProcessingResult {
	return NewSuccessResult()
}

type Processor func(ctx *ProcessingContext) ProcessingResult

func SuccessProcessor(ctx *ProcessingContext) ProcessingResult {
	return NewSuccessResult()
}

func Chain(stages ...Stage) Processor {

	next := SuccessStage

	wrap := func(stage Stage, processor Processor) Processor {
		return func(ctx *ProcessingContext) ProcessingResult {
			return stage(ctx, processor)
		}
	}

	for i := len(stages)-1; i >= 0; i-- {
		next = wrap(stages[i], next)
	}

	return next
}

func AsStage(processor Processor) Stage {
	return func(ctx *ProcessingContext, next Processor) ProcessingResult {
		result := processor(ctx)

		if result.ResultType != SuccessResult {
			return result
		}

		return next(ctx)
	}
}