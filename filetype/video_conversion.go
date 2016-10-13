package filetype

import (
	"github.com/jpg0/flickrup/processing"
	"path/filepath"
	"os/exec"
	"os"
	"github.com/juju/errors"
	"github.com/Sirupsen/logrus"
	"strings"
)

func VideoConversionStage() func(ctx *processing.PreprocessingContext, next processing.Preprocessor) processing.ProcessingResult {
	return func(ctx *processing.PreprocessingContext, next processing.Preprocessor) processing.ProcessingResult {

		ext := filepath.Ext(ctx.Filepath)

		if len(ext) > 0 {
			ext = ext[1:]
		}

		conversionCmd := ctx.Config.ConvertFiles[strings.ToLower(ext)]

		if conversionCmd != nil {
			out, err := convert(conversionCmd, ctx.Filepath)

			if err == nil {
				return processing.NewSuccessResult()
			} else {
				logrus.Warnf("Failed to convert video file %s: %s", ctx.Filepath, out)
				return processing.NewErrorResult(errors.Annotate(err, "Converting video"))
			}
		}

		return next(ctx)
	}
}

func convert(conversionCmd []string, filepath string) (string, error) {

	expand := func(k string) string {
		if k != "file" {
			panic("Only ${file} expansion supported")
		}

		return filepath
	}

	formatted := make([]string, len(conversionCmd) - 1)

	for i := range conversionCmd[1:] {
		formatted[i] = os.Expand(conversionCmd[i + 1], expand)
	}

	cmdOut, err := exec.Command(conversionCmd[0], formatted...).CombinedOutput()
	return string(cmdOut), err
}