package processing

import (
	"github.com/jpg0/flickrup/config"
	"time"
)

type ProcessingContext struct {
	File TaggedFile
	Visibilty string
	Config *config.Config
	ArchiveSubdir string
	UploadedId string
	OverrideDateTaken time.Time
	ArchivedAs string
}


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

func (pc ProcessingContext) DateTakenForArchive() time.Time {
	if pc.OverrideDateTaken.IsZero() {
		return pc.File.DateTaken()
	} else {
		return pc.OverrideDateTaken
	}
}

func NewProcessingContext(config *config.Config, file TaggedFile) *ProcessingContext {
	return &ProcessingContext{
		Visibilty: "public",
		Config: config,
		File: file,
	}
}