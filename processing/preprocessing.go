package processing

import (
	"github.com/jpg0/flickrup/config"
)

type PreprocessingContext struct {
	Filepath string
	Config *config.Config
	RequiresRestart bool
	changeSink ChangeSink
}

func NewPreprocessingContext(config *config.Config, filepath string, changeSink ChangeSink) *PreprocessingContext {
	return &PreprocessingContext{
		Config: config,
		Filepath: filepath,
		RequiresRestart: false,
		changeSink: changeSink,
	}
}

func (ppc *PreprocessingContext) ExpectChange() {
	ppc.changeSink.Expect(ppc.Filepath)
}