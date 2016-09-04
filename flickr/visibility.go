package flickr

import (
	"gopkg.in/masci/flickr.v2"
	"golang.org/x/net/context"
	log "github.com/Sirupsen/logrus"
)

/*
Returns whether file should be uploaded at all
 */
func AddVisibility(uploadParams *flickr.UploadParams, ctx context.Context) bool {

	visibility := ctx.Value("visibility").(string)
	switch visibility {
	default:
		panic("Unknown visibility: " + visibility)
	case "offline":
		log.Debug("Offline Visibility specified. Disabling upload.")
		return false
	case "family":
		log.Debug("Family visibility specified.")
		uploadParams.IsFamily = true
		uploadParams.IsFriend = false
		uploadParams.IsPublic = false
		return true
	case "friends":
		log.Debug("Friends visibility specified.")
		uploadParams.IsFamily = true
		uploadParams.IsFriend = true
		uploadParams.IsPublic = false
		return true
	case "private":
		log.Debug("Private visibility specified.")
		uploadParams.IsFamily = false
		uploadParams.IsFriend = false
		uploadParams.IsPublic = false
		return true
	case "public":
		log.Debug("Public visibility specified.")
		uploadParams.IsFamily = true
		uploadParams.IsFriend = true
		uploadParams.IsPublic = true
		return true
	}
}