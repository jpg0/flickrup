package filetype

import (
	"os"
	"github.com/Sirupsen/logrus"
	"io"
	"github.com/juju/errors"
)

func MoveFile(from string, to string) error {
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
