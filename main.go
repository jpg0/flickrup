package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/jpg0/flickrup/flickr"
	"github.com/jpg0/flickrup/filetype"
	flickrupconfig "github.com/jpg0/flickrup/config"

	"fmt"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/archive"
	"github.com/jpg0/flickrup/tags"
	"github.com/jpg0/flickrup/listen"
	"github.com/juju/errors"
	log "github.com/Sirupsen/logrus"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "flickrup"
	app.Usage = "Upload photos to Flickr"

	app.Commands = []cli.Command{
		cli.Command{
			Name: "upload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "config",
					Usage: "File path to configuration file",
				},
				cli.StringFlag{
					Name: "visibility",
					Usage: "public|friends|family|private|offline",
				},
			},
			Action: uploadFile,
		},
		cli.Command{
			Name: "watch",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "config",
					Usage: "File path to configuration file",
				},
				cli.StringFlag{
					Name: "loglevel",
					Usage: "Logging level",
				},
			},
			Action: verbose(watch),
		},
		cli.Command{
			Name: "test",
			Action: test,
		},
	}

	app.Run(os.Args)
}

func verbose(next func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		err := next(c)

		if err != nil {
			fmt.Println(errors.ErrorStack(err))
		}

		return err
	}
}

func initLogging(level string) error {
	switch strings.ToLower(level) {
	case "debug": log.SetLevel(log.DebugLevel)
	case "info": log.SetLevel(log.InfoLevel)
	case "warn": log.SetLevel(log.WarnLevel)
	case "error": log.SetLevel(log.ErrorLevel)
	case "fatal": log.SetLevel(log.FatalLevel)
	default:
		return errors.Errorf("Unknown logging level: %v", level)
	}

	return nil
}

func watch(c *cli.Context) error {

	err := initLogging(c.String("loglevel"))

	if err != nil {
		return errors.Trace(err)
	}

	config, err := flickrupconfig.Load(c.String("config"))

	if err != nil {
		return errors.Trace(err)
	}

	return listen.Watch(config)
}

func uploadFile(c *cli.Context) error {

	config, err := flickrupconfig.Load(c.String("config"))

	if err != nil {
		return err
	}

	taggedFile, err := filetype.NewTaggedImage(c.Args().First())

	if err != nil {
		return err
	}

	ctx := processing.NewProcessingContext()

	ctx.Visibilty = c.String("visibility")

	ctx.File = taggedFile

	ctx.Config = config

	processor, err := ProcessorPipeline(config)

	if err != nil {
		return err
	}

	err = processor(ctx)

	if err != nil {
		return err
	}

	return nil
}

func ProcessorPipeline(config *flickrupconfig.Config) (processing.Processor, error) {

	client, err := flickr.NewUploadClient(config)

	if err != nil {
		return nil, errors.Trace(err)
	}

	tagSetProcessor, err := tags.NewTagSetProcessor(config)
	rewriter := tags.NewRewriter()

	return processing.Chain(
		processing.AsStage(rewriter.MaybeRewrite),
		processing.AsStage(tags.MaybeBlock),
		processing.AsStage(tags.MaybeReplace),
		tagSetProcessor.Stage(),
		client.Stage(),
		processing.AsStage(archive.Archive),
	), nil
}

func test(c *cli.Context) error {
	fmt.Println(os.ExpandEnv("${HOME}/flic"))
	return nil
}