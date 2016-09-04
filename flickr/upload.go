package flickr

import (
	"gopkg.in/masci/flickr.v2"
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"bufio"
	"strings"
	"github.com/jpg0/flickrup/filetype"
	"golang.org/x/net/context"
	"errors"
	log "github.com/Sirupsen/logrus"

)

type UploadClient interface {
	Upload(file filetype.TaggedFile, ctx context.Context) error
}

type FlickrUploadClient struct {
	client *flickr.FlickrClient
}

func (client *FlickrUploadClient) Upload(file filetype.TaggedFile, ctx context.Context) error {

	params := flickr.NewUploadParams()

	if !AddVisibility(params, ctx) {
		//offline
		return nil
	}

	params.Tags = file.Keywords()

	response, err := flickr.UploadFile(client.client, file.File().Name(), params)

	if err != nil {
		return err
	}

	if response.HasErrors() {
		log.Errorf("Failed to upload photo %v: %v", file.File().Name(), response)
		return errors.New(response.ErrorMsg())
	} else {
		log.Debugf("Uploaded photo %v %v as %v", file.File().Name(), file.Keywords(), response.ID)
	}

	return nil
}

func NewClient(APIKey string, SharedSecret string) (*FlickrUploadClient, error){
	client := flickr.NewFlickrClient(APIKey, SharedSecret)
	token, err := getToken(client)

	if err != nil {
		return nil, err
	}

	client.OAuthToken = token.OAuthToken
	client.OAuthTokenSecret = token.OAuthTokenSecret

	return &FlickrUploadClient{client: client}, nil
}

func getToken(client *flickr.FlickrClient) (*flickr.OAuthToken, error) {
	filepath := fmt.Sprint(os.ExpandEnv("HOME"), "/.flickrup")

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		token, err := requestToken(client)

		if err != nil {
			return nil, err
		}

		bytes, err := yaml.Marshal(token)

		err = ioutil.WriteFile(filepath, bytes, 0644)

		if err != nil {
			return nil, err
		}

		return token, nil

	} else {
		//load from file

		file, err := os.Open(filepath)
		defer file.Close()

		if err != nil {
			return nil, err
		}

		tokens := new(flickr.OAuthToken)

		bytes, err := ioutil.ReadAll(file)

		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(bytes, tokens)

		if err != nil {
			return nil, err
		}

		return tokens, nil
	}
}

func requestToken(client *flickr.FlickrClient) (*flickr.OAuthToken, error){
	requestToken, err := flickr.GetRequestToken(client)

	if err != nil {
		return nil, err
	}

	url, err := flickr.GetAuthorizeUrl(client, requestToken)

	if err != nil {
		return nil, err
	}

	fmt.Println("Open this url in your process to complete the authentication process: ", url)
	fmt.Println("Copy here the number given when you complete the process.")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return flickr.GetAccessToken(client, requestToken, strings.TrimSpace(text))

}