package listen

import (
	"os"
	"io/ioutil"
	"strings"
	"github.com/kennygrant/sanitize"
	"fmt"
	"github.com/Sirupsen/logrus"
)

const STATUS_FILE_PREFIX = "_flickrup."

type UploadStatus struct {
	dir string
}

func NewUploadStatus(dir string) *UploadStatus {
	return &UploadStatus{dir:dir}
}

func (us *UploadStatus) IsStatusFile(path string) bool {
	return strings.HasPrefix(path, fmt.Sprintf("%v%v%v", us.dir, string(os.PathSeparator), STATUS_FILE_PREFIX))
}

func (us *UploadStatus) UpdateStatus(message string) (err error) {
	err = us.ClearStatus()

	if err != nil {
		return
	}

	return us.WriteStatus(message)
}

func (us *UploadStatus) WriteStatus(message string) (err error) {
	file, err := os.OpenFile(fmt.Sprintf("%v%v%v%v", us.dir, string(os.PathSeparator), STATUS_FILE_PREFIX, sanitize.BaseName(message)), os.O_RDONLY|os.O_CREATE, 0666)
	logrus.Debugf("Created file %v", file.Name())
	defer file.Close()
	return
}

func (us *UploadStatus) ClearStatus() (err error) {
	files, err := ioutil.ReadDir(us.dir)

	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), STATUS_FILE_PREFIX) {
			err = os.Remove(fmt.Sprintf("%v%v%v", us.dir, string(os.PathSeparator), file.Name()))
			if err != nil {
				logrus.Warnf("Failed to remove file %v", file.Name())
			} else {
				logrus.Debugf("Removed file %v", file.Name())
			}
		}
	}

	return
}