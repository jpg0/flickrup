package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/juju/errors"
	"github.com/Sirupsen/logrus"
	flickrupconfig "github.com/jpg0/flickrup/config"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "flickrup"
	app.Usage = "Upload photos to Flickr"
	app.Version = "1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config",
			Usage: "File path to configuration file",
		},
		cli.StringFlag{
			Name: "loglevel",
			Usage: "Logging level",
			Value: "info",
		},
	}
	app.Action = verbose(watch)
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
	case "debug": logrus.SetLevel(logrus.DebugLevel)
	case "info": logrus.SetLevel(logrus.InfoLevel)
	case "warn": logrus.SetLevel(logrus.WarnLevel)
	case "error": logrus.SetLevel(logrus.ErrorLevel)
	case "fatal": logrus.SetLevel(logrus.FatalLevel)
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