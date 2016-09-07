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

type TransferResponse struct {
	Id string
}

func Transfer(filepath string, tags []string, isPublic bool, isFamily bool, isFriend bool, servicePassword string) (string, error) {

	file, err := os.Open(filepath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	idx := strings.Index(filepath, "/Camera Uploads/")

	if idx == -1 {
		return "", errors.New("Cannot find /Camera Uploads/ in path")
	}

	title := filepath[idx:]

	stat, err := file.Stat()

	if err != nil {
		return "", err
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
		return "", err
	}

	req, err := http.NewRequest("POST", "http://d2f-transfer.appspot.com/transfer", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if servicePassword != "" {
		req.SetBasicAuth("flickrup", servicePassword)

	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	resBody, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusCreated {
		if err != nil {
			return "", err
		}

		return "", errors.New(fmt.Sprintf("%v: %v", res.StatusCode, string(resBody)))
	}

	transferResponse := TransferResponse{}

	json.Unmarshal(resBody, transferResponse)

	return transferResponse.Id, nil
}