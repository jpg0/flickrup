package tags

import ("github.com/jpg0/flickrup/processing"
	log "github.com/Sirupsen/logrus"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/flickr"
	"github.com/juju/errors"
)

type TagSetProcessor struct {
	setClient flickr.SetClient
}

func NewTagSetProcessor(config *config.Config) (*TagSetProcessor, error) {

	setClient, err := flickr.NewFlickrSetClient(config.APIKey, config.SharedSecret)

	if err != nil {
		return nil, err
	}

	return &TagSetProcessor{
		setClient: setClient,
	}, nil
}

//must occur between upload & archive stages
func (tsp *TagSetProcessor) Stage() processing.Stage {
	return func(ctx *processing.ProcessingContext, next processing.Processor) error {

		sets := processing.ValuesByPrefix(ctx.File.Keywords(), ctx.Config.TagsetPrefix)

		if ctx.Visibilty != "offline" {
			if len(sets) > 0 {

				log.Debugf("Detected membership of set(s) %v", sets)


				for _, set := range sets {
					log.Debugf("Adding %v to set: %v", ctx.File.Name(), set)
					err := tsp.setClient.AddToSet(ctx.UploadedId, set, ctx.File.DateTaken())

					if err != nil {
						return errors.Annotate(err, "Adding photo to set")
					}
				}

				ctx.ArchiveSubdir = sets[0]
				date, err := tsp.setClient.DateOfSet(sets[0])

				if err != nil {
					return errors.Annotate(err, "Getting date of set")
				}

				ctx.OverrideDateTaken = date
			}
		}

		return next(ctx)
	}
}