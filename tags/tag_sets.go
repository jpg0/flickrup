package tags

import ("github.com/jpg0/flickrup/processing"
	log "github.com/Sirupsen/logrus"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/flickraccess"
	"github.com/juju/errors"
)

type TagSetProcessor struct {
	setClient flickraccess.SetClient
}

func NewTagSetProcessor(config *config.Config) (*TagSetProcessor, error) {

	setClient, err := flickraccess.NewFlickrSetClient(config.APIKey, config.SharedSecret)

	if err != nil {
		return nil, err
	}

	return &TagSetProcessor{
		setClient: setClient,
	}, nil
}

//must occur between upload & archive stages
func (tsp *TagSetProcessor) Stage() processing.Stage {
	return func(ctx *processing.ProcessingContext, next processing.Processor) processing.ProcessingResult {

		sets := processing.ValuesByPrefix(ctx.File.Keywords(), ctx.Config.TagsetPrefix)

		if ctx.Visibilty != "offline" {
			if len(sets) > 0 {

				log.Infof("Detected membership of set(s) %v", sets)


				for _, set := range sets {
					log.Infof("Adding %v to set: %v", ctx.File.Name(), set)
					err := tsp.setClient.AddToSet(ctx.UploadedId, set, ctx.File.DateTaken())

					if err != nil {
						return processing.NewErrorResult(errors.Annotate(err, "Adding photo to set"))
					}
				}

				ctx.ArchiveSubdir = sets[0]
				date, err := tsp.setClient.DateOfSet(sets[0])

				if err != nil {
					return processing.NewErrorResult(errors.Annotate(err, "Getting date of set"))
				}

				ctx.OverrideDateTaken = date
			}
		}

		return next(ctx)
	}
}