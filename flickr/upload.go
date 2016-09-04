package flickr

import (
	"gopkg.in/masci/flickr.v2"
	"github.com/jpg0/flickrup/filetype"
	flickrupconfig "github.com/jpg0/flickrup/config"
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

func (client *FlickrUploadClient) Upload(file filetype.TaggedFile, ctx *filetype.ProcessingContext, cfg *flickrupconfig.Config) error {

	params := flickr.NewUploadParams()

	if !AddVisibility(params, ctx) {
		//offline
		return nil
	}

	params.Tags = file.Keywords()

	if cfg.TransferServicePassword != "" {
		err := Transfer(file.Filepath(), params.Tags, params.IsPublic, params.IsFamily, params.IsFriend, cfg.TransferServicePassword)

		if err != nil {
			log.Infof("Failed to transfer: %v", err)
			log.Info("Falling back to direct upload")
		} else {
			return nil
		}
	}

	response, err := flickr.UploadFile(client.client, file.Filepath(), params)

	if err != nil {
		return err
	}

	if response.HasErrors() {
		log.Errorf("Failed to upload photo %v: %v", file.Name(), response)
		return errors.New(response.ErrorMsg())
	} else {
		log.Debugf("Uploaded photo %v %v as %v", file.Name(), file.Keywords(), response.ID)
	}

	return nil
}

func NewUploadClient(APIKey string, SharedSecret string) (*FlickrUploadClient, error){
	client := flickr.NewFlickrClient(APIKey, SharedSecret)
	token, err := getToken(client)

	if err != nil {
		return nil, err
	}

	client.OAuthToken = token.OAuthToken
	client.OAuthTokenSecret = token.OAuthTokenSecret

	return &FlickrUploadClient{client: client}, nil
}