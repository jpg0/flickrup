package processing

import (
	"github.com/jpg0/flickrup/config"
)

type PreprocessingContext struct {
	Filepath string
	Config *config.Config
	RequiresRestart bool
}

func NewPreprocessingContext(config *config.Config, filepath string) *PreprocessingContext {
	return &PreprocessingContext{
		Config: config,
		Filepath: filepath,
		RequiresRestart: false,
	}
}