package main

import (
	"os"
	"io/ioutil"
	"strings"
	"github.com/kennygrant/sanitize"
	"fmt"
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
	file, err := os.OpenFile(fmt.Sprintf("%s%s%s", dir, os.PathSeparator, sanitize.BaseName(message)), os.O_RDONLY|os.O_CREATE, 0666)
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
			defer os.Remove(fmt.Sprintf("%s%s%s", dir, os.PathSeparator, file.Name()))
		}
	}

	return
}