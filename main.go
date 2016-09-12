package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/juju/errors"
	log "github.com/Sirupsen/logrus"
	flickrupconfig "github.com/jpg0/flickrup/config"
	"strings"
	"github.com/jpg0/flickrup/flickr"
)

func main() {
	app := cli.NewApp()
	app.Name = "flickrup"
	app.Usage = "Upload photos to Flickr"

	app.Commands = []cli.Command{
		{
			Name: "watch",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "config",
					Usage: "File path to configuration file",
				},
				cli.StringFlag{
					Name: "loglevel",
					Usage: "Logging level",
					Value: "info",
				},
			},
			Action: verbose(watch),
		},
		{
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

	return CreateAndRunPipeline(config)
}

func test(c *cli.Context) error {
	flickr.DT()
	return nil
}