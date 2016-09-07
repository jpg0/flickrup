package flickr

import (
	"gopkg.in/masci/flickr.v2"
	"github.com/jpg0/flickrup/processing"
	"golang.org/x/net/context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/jpg0/flickrup/config"
)

type UploadClient interface {
	Upload(file processing.TaggedFile, ctx context.Context) error
}

type FlickrUploadClient struct {
	client *flickr.FlickrClient
}

func (client *FlickrUploadClient) Stage() processing.Stage {
	return func(ctx *processing.ProcessingContext, next processing.Processor) error {
		err := client.Upload(ctx)

		if err != nil {
			return err
		}

		return next(ctx)
	}
}

func (client *FlickrUploadClient) Upload(ctx *processing.ProcessingContext) error {

	params := flickr.NewUploadParams()

	if !AddVisibility(params, ctx) {
		//offline
		return nil
	}

	file := ctx.File

	params.Tags = file.Keywords().All()

	if ctx.Config.TransferServicePassword != "" {
		id, err := Transfer(file.Filepath(), params.Tags, params.IsPublic, params.IsFamily, params.IsFriend, ctx.Config.TransferServicePassword)

		if err != nil {
			log.Infof("Failed to transfer: %v", err)
			log.Info("Falling back to direct upload")
		} else {
			ctx.UploadedId = id
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
		ctx.UploadedId = response.ID
	}

	return nil
}

func NewUploadClient(config *config.Config) (*FlickrUploadClient, error){
	client := flickr.NewFlickrClient(config.APIKey, config.SharedSecret)
	token, err := getToken(client)

	if err != nil {
		return nil, err
	}

	client.OAuthToken = token.OAuthToken
	client.OAuthTokenSecret = token.OAuthTokenSecret

	return &FlickrUploadClient{client: client}, nil
}