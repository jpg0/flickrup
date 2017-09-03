package flickraccess

import (
	"time"
	"errors"
	"github.com/jpg0/flickr"
	"github.com/jpg0/flickr/photos"
	"github.com/jpg0/flickr/photosets"
)

type SetClient interface {
	DateOfSet(setName string) (time.Time, error)
	AddToSet(photoId string, setName string, datePhotoTaken time.Time) error
}

type FlickrSetClient struct {
	flickrClient *flickr.FlickrClient
	setIdToDate map[string]time.Time
	setNameToId map[string]string
}

func NewFlickrSetClient(APIKey string, SharedSecret string) (SetClient, error){
	client := flickr.NewFlickrClient(APIKey, SharedSecret)
	token, err := getToken(client)

	if err != nil {
		return nil, err
	}

	client.OAuthToken = token.OAuthToken
	client.OAuthTokenSecret = token.OAuthTokenSecret

	return &FlickrSetClient{
		flickrClient: client,
		setIdToDate: make(map[string]time.Time),
		setNameToId: make(map[string]string),
	}, nil
}

func (client FlickrSetClient) DateOfSet(setName string) (time.Time, error) {
	id, err := client.setIdFromName(setName)

	if err != nil {
		return time.Time{}, err
	}

	date := client.setIdToDate[id]

	if date.IsZero() {
		setResponse, err := photosets.GetInfo(client.flickrClient, true, id, "")

		if err != nil {
			return time.Time{}, err
		}

		if setResponse.HasErrors() {
			return time.Time{}, errors.New(setResponse.ErrorMsg())
		}

		primary := setResponse.Set.Primary

		photoResponse, err := photos.GetInfo(client.flickrClient, primary, "")

		if err != nil {
			return time.Time{}, err
		}

		if photoResponse.HasErrors() {
			return time.Time{}, errors.New(photoResponse.ErrorMsg())
		}

		timeAsString := photoResponse.Photo.Dates.Taken
		date, err = time.Parse(time.RFC3339, timeAsString)

		if err != nil {
			return time.Time{}, err
		}

		client.setIdToDate[id] = date
	}

	return date, nil
}

func (client FlickrSetClient) setIdFromName(setName string) (string, error) {
	val := client.setNameToId[setName]

	if val == "" {
		response, err := photosets.GetList(client.flickrClient, true, "", 0)

		if err != nil {
			return "", err
		}

		if response.HasErrors() {
			return "", errors.New(response.ErrorMsg())
		}

		for _, set := range response.Photosets.Items {
			client.setNameToId[set.Title] = set.Id
		}

		val = client.setNameToId[setName]
	}

	return val, nil
}

func (client FlickrSetClient) AddToSet(photoId string, setName string, datePhotoTaken time.Time) error {
	setId, err := client.setIdFromName(setName)

	if err != nil {
		return err
	}

	if setId != "" {
		response, err := photosets.AddPhoto(client.flickrClient, setId, photoId)

		if err != nil {
			return err
		}

		if response.HasErrors() {
			return errors.New(response.ErrorMsg())
		}

	} else { //if we still don't, create it

		response, err := photosets.Create(client.flickrClient, setName, "", photoId)

		if err != nil {
			return err
		}

		if response.HasErrors() {
			return errors.New(response.ErrorMsg())
		}

		client.setNameToId[setName] = response.Set.Id
		client.setIdToDate[response.Set.Id] = datePhotoTaken
	}

	return nil
}