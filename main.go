package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/jpg0/flickrup/flickr"
)

func main() {
	app := cli.NewApp()
	app.Name = "flickrup"
	app.Usage = "Upload photos to Flickr"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config",
			Usage: "File path to configuration file",
		},
		cli.BoolFlag{
			Name: "friends",
			Usage: "Visible to friends only?",
		},
		cli.BoolFlag{
			Name: "family",
			Usage: "Visible to family only?",
		},
		cli.BoolFlag{
			Name: "private",
			Usage: "Visible to me only?",
		},
		cli.StringSliceFlag{
			Name: "tags",
			Usage: "tags for file",
		},
		cli.StringFlag{
			Name: "password",
			Usage: "password to authorise transfer",
		},
	}

	app.Action = run
	app.Run(os.Args)
}


func run(c *cli.Context) error {

	file, err := os.Open(c.Args().First())

	if err != nil {
		return err
	}

	defer file.Close()

	return flickr.Transfer(
		file,
		c.StringSlice("tags"),
		c.Bool("public"),
		c.Bool("family"),
		c.Bool("friend"),
		c.String("password"),
	)
}