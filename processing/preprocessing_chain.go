package processing

type PreStage func(ctx *PreprocessingContext, next Preprocessor) ProcessingResult

func SuccessPreStage(_ *PreprocessingContext) ProcessingResult {
	return NewSuccessResult()
}

type Preprocessor func(ctx *PreprocessingContext) ProcessingResult

func SuccessPreprocessor(ctx *PreprocessingContext) ProcessingResult {
	return NewSuccessResult()
}

func ChainPreStages(stages ...PreStage) Preprocessor {

	next := SuccessPreStage

	wrap := func(stage PreStage, preprocessor Preprocessor) Preprocessor {
		return func(ctx *PreprocessingContext) ProcessingResult {
			return stage(ctx, preprocessor)
		}
	}

	for i := len(stages)-1; i >= 0; i-- {
		next = wrap(stages[i], next)
	}

	return next
}