package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
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

func clearCache(apiKey string, email string, files []string) error {
	for i := range files {
		files[i] = "https://assets.honestbee.com/" + files[i]
	}
	purgeCacheRequest := cloudflare.PurgeCacheRequest{
		Files: files,
	}

	// Construct a new API object
	api, err := cloudflare.New(apiKey, email)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName("honestbee.com") // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	purgeCacheRespone, err := api.PurgeCache(id, purgeCacheRequest)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Printf("status: %v", purgeCacheRespone.Response.Success)
	return err
}

func parseFile(fileName string) ([]string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var fileSlice stringsAlias
	fileSlice = strings.Split(string(file), "\n")
	return fileSlice.removeEmpty(), nil
}

func (s stringsAlias) removeEmpty() []string {
	for i, element := range s {
		if element == "" {
			s = append(s[:i], s[i+1:]...)
		}
	}
	return s
}
