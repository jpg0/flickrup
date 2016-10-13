package archive

import (
	"time"
	"fmt"
	"os"
	"path/filepath"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
	"github.com/jpg0/flickrup/filetype"
)

func Archive(ctx *processing.ProcessingContext) processing.ProcessingResult {
	newPath, err := archiveFileByDate(ctx.File.Filepath(), ctx.Config.ArchiveDir, ctx.DateTakenForArchive(), ctx.ArchiveSubdir)

	if err != nil {
		return processing.NewErrorResult(errors.Trace(err))
	}

	ctx.ArchivedAs = newPath

	return processing.NewSuccessResult()
}

func archiveFileByDate(file string, toDir string, date time.Time, subdir string) (string, error) {
	targetDir := fmt.Sprintf("%v/%v/%.2d/%v", toDir, date.Year(), date.Month(), subdir)
	err := os.MkdirAll(targetDir, 0755)

	if err != nil {
		return "", errors.Trace(err)
	}

	newName := fmt.Sprintf("%v/%v", targetDir, filepath.Base(file))

	err = filetype.MoveFile(file, newName)

	if err != nil {
		return "", errors.Trace(err)
	}

	return newName, nil
}