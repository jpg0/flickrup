package archive

import (
	"time"
	"fmt"
	"os"
	"path/filepath"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
)

func Archive(ctx *processing.ProcessingContext) error {
	return archiveFileByDate(ctx.File.Filepath(), ctx.Config.ArchiveDir, ctx.DateTaken(), ctx.ArchiveSubdir)
}

func archiveFileByDate(file string, toDir string, date time.Time, subdir string) error {
	targetDir := fmt.Sprintf("%v/%v/%.2v/%v", toDir, date.Year(), date.Month(), subdir)
	err := os.MkdirAll(targetDir, 0755)

	if err != nil {
		return errors.Trace(err)
	}

	newName := fmt.Sprintf("%v/%v", targetDir, filepath.Base(file))

	return os.Rename(file, newName)
}