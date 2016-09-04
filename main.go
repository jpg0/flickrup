package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/jpg0/flickrup/flickr"
	"github.com/jpg0/flickrup/filetype"
	flickrupconfig "github.com/jpg0/flickrup/config"

	"fmt"
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
			Name: "test",
			Action: test,
		},
	}

	app.Run(os.Args)
}


func uploadFile(c *cli.Context) error {

	config, err := flickrupconfig.Load(c.String("config"))

	if err != nil {
		return err
	}

	client, err := flickr.NewUploadClient(config.APIKey, config.SharedSecret)

	if err != nil {
		return err
	}

	taggedFile, err := filetype.NewTaggedImage(c.Args().First())

	if err != nil {
		return err
	}

	ctx := filetype.NewProcessingContext()

	ctx.Visibilty = c.String("visibility")

	err = client.Upload(taggedFile, ctx, config)

	if err != nil {
		return err
	}

	return nil
}

func test(c *cli.Context) error {
	fmt.Println(os.ExpandEnv("${HOME}/flic"))
	return nil
}