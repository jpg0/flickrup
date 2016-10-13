package processing

type ProcessingResultType int

const (
	SuccessResult ProcessingResultType = iota
	ErrorResult ProcessingResultType = iota
	RestartResult ProcessingResultType = iota
)

type ProcessingResult struct {
	ResultType ProcessingResultType
	Error error
}

func NewSuccessResult() ProcessingResult {
	return ProcessingResult{
		ResultType: SuccessResult,
	}
}

func NewRestartResult() ProcessingResult {
	return ProcessingResult{
		ResultType: RestartResult,
	}
}

func NewErrorResult(err error) ProcessingResult {
	return ProcessingResult{
		ResultType: ErrorResult,
		Error: err,
	}
}