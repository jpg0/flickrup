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
	FileUpdated bool
	changeSink ChangeSink
}

func (pc ProcessingContext) DateTakenForArchive() time.Time {
	if pc.OverrideDateTaken.IsZero() {
		return pc.File.DateTaken()
	} else {
		return pc.OverrideDateTaken
	}
}

func NewProcessingContext(config *config.Config, file TaggedFile, changeSink ChangeSink) *ProcessingContext {
	return &ProcessingContext{
		Visibilty: "default",
		Config: config,
		File: file,
		FileUpdated: false,
		changeSink: changeSink,
	}
}

func (pc *ProcessingContext) ExpectChange() {
	pc.changeSink.Expect(pc.File.Filepath())
	pc.FileUpdated = true
}