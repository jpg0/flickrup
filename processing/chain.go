package processing

type Stage func(ctx *ProcessingContext, next Processor) error

func SuccessStage(_ *ProcessingContext) error {
	return nil
}

type Processor func(ctx *ProcessingContext) error

func SuccessProcessor(ctx *ProcessingContext) error {
	return nil
}

func Chain(stages ...Stage) Processor {

	next := SuccessStage

	wrap := func(stage Stage, processor Processor) Processor {
		return func(ctx *ProcessingContext) error {
			return stage(ctx, processor)
		}
	}

	for i := len(stages)-1; i >= 0; i-- {
		next = wrap(stages[i], next)
	}

	return next
}

func AsStage(processor Processor) Stage {
	return func(ctx *ProcessingContext, next Processor) error {
		err := processor(ctx)

		if err != nil {
			return err
		}

		return next(ctx)
	}
}