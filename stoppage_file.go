package main

import (
	"os"
	"io/ioutil"
	"strings"
	"github.com/kennygrant/sanitize"
	"fmt"
	"github.com/Sirupsen/logrus"
)

const STATUS_FILE_PREFIX = "_flickrup."

func UpdateStatus(message string, dir string) (err error) {
	err = ClearStatus(dir)

	if err != nil {
		return
	}

	return WriteStatus(message, dir)
}

func WriteStatus(message string, dir string) (err error) {
	file, err := os.OpenFile(fmt.Sprintf("%s%s%s%s", dir, os.PathSeparator, STATUS_FILE_PREFIX, sanitize.BaseName(message)), os.O_RDONLY|os.O_CREATE, 0666)
	logrus.Debugf("Created file %s", file.Name())
	defer file.Close()
	return
}

func ClearStatus(dir string) (err error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), STATUS_FILE_PREFIX) {
			err = os.Remove(fmt.Sprintf("%s%s%s", dir, os.PathSeparator, file.Name()))
			if err != nil {
				logrus.Warnf("Failed to remove file %s", file.Name())
			} else {
				logrus.Debugf("Removed file %s", file.Name())
			}
		}
	}

	return
}