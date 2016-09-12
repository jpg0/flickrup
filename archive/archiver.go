package archive

import (
	"time"
	"fmt"
	"os"
	"path/filepath"
	"github.com/jpg0/flickrup/processing"
	"github.com/juju/errors"
	"github.com/Sirupsen/logrus"
	"io"
)

func Archive(ctx *processing.ProcessingContext) error {
	return archiveFileByDate(ctx.File.Filepath(), ctx.Config.ArchiveDir, ctx.DateTakenForArchive(), ctx.ArchiveSubdir)
}

func archiveFileByDate(file string, toDir string, date time.Time, subdir string) error {
	targetDir := fmt.Sprintf("%v/%v/%.2d/%v", toDir, date.Year(), date.Month(), subdir)
	err := os.MkdirAll(targetDir, 0755)

	if err != nil {
		return errors.Trace(err)
	}

	newName := fmt.Sprintf("%v/%v", targetDir, filepath.Base(file))

	return moveFile(file, newName)
}

func moveFile(from string, to string) error {
	err := os.Rename(from, to)

	if err == nil {
		return nil
	}

	logrus.Debugf("Failed to move file: %v", err)
	logrus.Debugf("Falling back to copy...")

	err = copyFileContents(from, to)

	if err != nil {
		return errors.Annotate(err, "Failed to copy file for archive")
	}

	return os.Remove(from)
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return nil
}