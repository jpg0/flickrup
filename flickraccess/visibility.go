package flickraccess

import (
	"github.com/jpg0/flickr"
	"github.com/jpg0/flickrup/processing"
)

/*
Returns whether file should be uploaded at all
 */
func AddVisibility(uploadParams *flickr.UploadParams, ctx *processing.ProcessingContext) bool {

	visibility := ctx.Visibilty
	switch visibility {
	default:
		panic("Unknown visibility: " + visibility)
	case "default":
		uploadParams.IsDefault = true
		return true
	case "offline":
		return false
	case "family":
		uploadParams.IsFamily = true
		uploadParams.IsFriend = false
		uploadParams.IsPublic = false
		return true
	case "friends":
		uploadParams.IsFamily = true
		uploadParams.IsFriend = true
		uploadParams.IsPublic = false
		return true
	case "private":
		uploadParams.IsFamily = false
		uploadParams.IsFriend = false
		uploadParams.IsPublic = false
		return true
	case "public":
		uploadParams.IsFamily = true
		uploadParams.IsFriend = true
		uploadParams.IsPublic = true
		return true
	}
}