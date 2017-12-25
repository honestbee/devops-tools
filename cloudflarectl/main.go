package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

type stringsAlias []string

func main() {
	app := cli.NewApp()
	app.Name = "cloudflarectl"
	app.Usage = "golang tool to clear cloudflare cache"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Tuan Nguyen",
			Email: "tuan.nguyen@honestbee.com",
		},
	}
	app.Version = "0.9.0"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "apiKey",
			Usage:  "Cloudflare's API key (REQUIRED)",
			EnvVar: "CF_API_KEY",
		},
		cli.StringFlag{
			Name:   "email",
			Usage:  "Cloudflare's email account (REQUIRED)",
			EnvVar: "CF_API_EMAIL",
		},
		cli.StringFlag{
			Name:   "file, f",
			Value:  "./files_list.txt",
			Usage:  "`<files slice>` which need to be cleared ",
			EnvVar: "CF_FILES",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "clear",
			Aliases: []string{"c"},
			Usage:   "Clear list of files's cache",
			Flags:   flags,
			Action: func(c *cli.Context) error {
				apiKey := c.String("apiKey")
				email := c.String("email")
				file := c.String("file")

				if (len(email) == 0) || (len(apiKey) == 0) {
					return cli.NewExitError("CF_API_KEY & CF_API_EMAIL must be defined!", 1)
				}

				fileList, err := parseFile(file)
				if err != nil {
					log.Fatal(err)
				}

				err = clearCache(apiKey, email, fileList)
				return err
			},
		},
		{
			Name:    "clearAll",
			Aliases: []string{"ca"},
			Usage:   "Clear everything",
			Flags:   flags,
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Show account status",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)

}
