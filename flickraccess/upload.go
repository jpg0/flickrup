package flickraccess

import (
	"github.com/jpg0/flickr"
	"github.com/jpg0/flickrup/processing"
	"golang.org/x/net/context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/jpg0/flickrup/config"
)

type UploadClient interface {
	Upload(file processing.TaggedFile, ctx context.Context) processing.ProcessingResult
}

type FlickrUploadClient struct {
	client *flickr.FlickrClient
}

func (client *FlickrUploadClient) Stage() processing.Stage {
	return func(ctx *processing.ProcessingContext, next processing.Processor) processing.ProcessingResult {
		result := client.Upload(ctx)

		if result.ResultType != processing.SuccessResult {
			return result
		}

		return next(ctx)
	}
}

func (client *FlickrUploadClient) Upload(ctx *processing.ProcessingContext) processing.ProcessingResult {

	params := flickr.NewUploadParams()

	if !AddVisibility(params, ctx) {
		//offline
		log.Infof("Offline visibility specified; skipping upload")
		return processing.NewSuccessResult()
	}

	file := ctx.File

	params.Tags = file.Keywords().All().Slice()

	if ctx.Config.TransferService != nil {
		id, err := Transfer(ctx.Config.TransferService.MapDropboxPath(file.Filepath()), params.Tags, params.IsDefault, params.IsPublic, params.IsFamily, params.IsFriend, ctx.FileUpdated, ctx.Config.TransferService.Password)

		if err != nil {
			log.Infof("Failed to transfer: %v", err)
			log.Info("Falling back to direct upload")
		} else {
			ctx.UploadedId = id
			return processing.NewSuccessResult()
		}
	}

	response, err := flickr.UploadFile(client.client, file.Filepath(), params)

	if err != nil {
		return processing.NewErrorResult(err)
	}

	if response.HasErrors() {
		log.Errorf("Failed to upload photo %v: %v", file.Name(), response)
		return processing.NewErrorResult(errors.New(response.ErrorMsg()))
	} else {
		log.Infof("Uploaded photo %v %v as %v", file.Name(), file.Keywords().All().Slice(), response.ID)
		ctx.UploadedId = response.ID
	}

	return processing.NewSuccessResult()
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