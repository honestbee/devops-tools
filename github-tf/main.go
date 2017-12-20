package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "org, o",
			Usage:  "`organization` to generate tf config for",
			EnvVar: "GITHUB_ORGANIZATION",
		},
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "`token` to access GitHub API",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "error",
			Usage:  "Log level (panic, fatal, error, warn, info, or debug)",
			EnvVar: "LOG_LEVEL",
		},
	}
	app := cli.NewApp()
	app.Name = "github-tf"
	app.Usage = "Download GitHub teams to TF config"
	app.Action = run

	app.Version = fmt.Sprintf("0.1.%s", build)
	app.Author = "Honestbee DevOps"

	app.Flags = flags

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	logLevelString := c.String("log-level")
	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	gitHub := GitHub{
		Organization: c.String("org"),
		Token:        c.String("token"),
	}

	if gitHub.Organization == "" || gitHub.Token == "" {
		cli.ShowAppHelpAndExit(c, 1)
	}

	//gitHub.ListRepos()

	teams, err := gitHub.ListTeams()
	if err != nil {
		return err
	}

	for _, t := range teams {
		teamRoles, err := gitHub.GetTeamRoles(t)
		if err != nil {
			return err
		}
		// render template to a TF config per team
		f, err := os.Create(fmt.Sprintf("output/teams-config/%v.tf", *t.Slug))
		if err != nil {
			return err
		}
		err = RenderTerraformConfig(teamRoles, f)
		f.Close()

		// render template for an import bash script per team
		f, err = os.Create(fmt.Sprintf("output/teams-import/%v.sh", *t.Slug))
		if err != nil {
			return err
		}
		err = RenderTerraformImport(teamRoles, f)
		f.Close() //file won't be closed on panic? (use anynomous func and defer?)

		tempString := strings.Split(*t.Slug, "-")
		if len(tempString) == 2 {
			f, err = os.Create(fmt.Sprintf("output/vault-policies/%v.hcl", tempString))
			if err != nil {
				return err
			}
			err = RenderVaultPolicy(tempString, f)
			f.Close()
		}
	}
	return nil
}
