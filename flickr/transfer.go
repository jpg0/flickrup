package flickr

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/jpg0/dropbox"
)

type Validations struct {
	Mtime *dropbox.DBTime
	Size  *int64
}

type TransferRequest struct {
	Title                        string
	Tags                         []string
	IsPublic, IsFamily, IsFriend bool
	Validations                  *Validations
}

func Transfer(file *os.File, tags []string, isPublic bool, isFamily bool, isFriend bool, servicePassword string) error {

	title := file.Name()[strings.Index(file.Name(), "/Camera Uploads/"):]

	stat, err := file.Stat()

	if err != nil {
		return err
	}

	mtime := dropbox.DBTime(stat.ModTime())
	size := stat.Size()

	transferRequest := TransferRequest{
		Title: title,
		IsPublic: isPublic,
		IsFamily: isFamily,
		IsFriend: isFriend,
		Tags: tags,
		Validations: &Validations{
			Mtime: &mtime,
			Size: &size,
		},
	}

	body, err := json.Marshal(transferRequest)

	if err != nil {
		return err
	}

	//fmt.Print(string(body))

	req, err := http.NewRequest("POST", "http://d2f-transfer.appspot.com/transfer", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if servicePassword != "" {
		req.SetBasicAuth("flickrup", servicePassword)

	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("%v: %v", res.StatusCode, string(body)))
	}

	return nil
}