package flickr

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	"github.com/jpg0/dropbox"
	"github.com/juju/errors"
	"github.com/jpg0/flickrup/config"
	"gopkg.in/yaml.v2"
	"github.com/Sirupsen/logrus"
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

type DropboxConfig struct {
	Personal DropboxAccount `json:"personal"`
}

type DropboxAccount struct {
	Path string `json:"path"`
	Host uint64 `json:"host"`
	IsTeam bool `json:"is_team"`
	SubscriptionType string `json:"subscription_type"`
}

func DropboxDir() (string, error) {
	data, err := ioutil.ReadFile(os.ExpandEnv("${HOME}/.dropbox/info.json"))

	if err != nil {
		return "", errors.Trace(err)
	}

	dbc := &DropboxConfig{}

	err = json.Unmarshal(data, dbc)

	if err != nil {
		return "", errors.Trace(err)
	}

	return dbc.Personal.Path, nil
}

func DT() {
	dbc := &config.Config{
		TransferService: &config.TransferService{
			Password:"test",
			DropboxDirMapping:map[string]string {"a": "b"},
		},
	}

	b, _ := yaml.Marshal(dbc)

	fmt.Println(string(b))
}

func Transfer(filepath string, tags []string, isPublic bool, isFamily bool, isFriend bool, servicePassword string) (string, error) {

	file, err := os.Open(filepath)

	if err != nil {
		return "", errors.Trace(err)
	}

	defer file.Close()

	root, err := DropboxDir()

	if err != nil {
		return "", errors.Trace(err)
	}

	if !strings.HasPrefix(filepath, root) {
		return "", errors.New("Cannot find dropbox dir in path")
	}

	title := filepath[len(root):]

	stat, err := file.Stat()

	if err != nil {
		return "", err
	}

	mtime := dropbox.DBTime(stat.ModTime())
	size := stat.Size()

	transferRequest := &TransferRequest{
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
		return "", errors.Trace(err)
	}

	req, err := http.NewRequest("POST", "http://d2f-transfer.appspot.com/transfer", bytes.NewReader(body))
	if err != nil {
		return "", errors.Trace(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if servicePassword != "" {
		req.SetBasicAuth("flickrup", servicePassword)

	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", errors.Trace(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusCreated {
		if err != nil {
			return "", errors.Trace(err)
		}

		return "", errors.New(fmt.Sprintf("%v: %v", res.StatusCode, string(resBody)))
	}

	transferResponse := &TransferResponse{}

	err = json.Unmarshal(resBody, transferResponse)

	if err != nil {
		logrus.Debugf("Response was: %v", string(resBody))
		return "", errors.Annotate(err, "Failed reading transfer response")
	}

	logrus.Debugf("Successfully transfered %v", title)

	return transferResponse.Id, nil
}