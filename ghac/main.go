package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "log-level",
			Value:  "error",
			Usage:  "Log level (panic, fatal, error, warn, info, or debug)",
			EnvVar: "LOG_LEVEL",
		},
		cli.StringFlag{
			Name:  "source, s",
			Value: "teams.yaml",
			Usage: "Source `yaml`",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "output/teams-config",
			Usage: "Destination `directory` to render in - must exist",
		},
		cli.StringFlag{
			Name:  "template, t",
			Value: "templates/team.tf.tpl",
			Usage: "Desired template used to render output",
		},
	}
	app := cli.NewApp()
	app.Name = "ghac"
	app.Usage = "Manage GitHub Teams and Team membership in yaml"
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

	// log.Debugf("source: %v", c.String("source"))
	// log.Debugf("destination: %v", c.String("destination"))
	// log.Debugf("template: %v", c.String("template"))

	//TODO: Add more validations
	dstDirName := path.Dir(c.String("destination"))
	log.Debugf("dstDirName: %v", dstDirName)
	if stat, err := os.Stat(dstDirName); err != nil || !stat.IsDir() {
		fmt.Printf("Invalid destination: %v\n", dstDirName)
		cli.ShowAppHelpAndExit(c, 1)
	}

	tl, err := readTeams(c.String("source"))
	if err != nil {
		return err
	}

	tpl := c.String("template")                                    //expect name.tf.tpl or name.hcl.tpl
	suffix := path.Ext(strings.TrimSuffix(tpl, filepath.Ext(tpl))) // .tf or .hcl
	log.Debugf("suffix: %v", suffix)

	for _, t := range tl.Teams {

		// Check format
		if suffix == ".tf" {
			// render template to destination (1 file per team)
			f, err := os.Create(path.Join(dstDirName, fmt.Sprintf("%v%v", t.Slug, suffix)))
			if err != nil {
				return err
			}

			err = RenderTemplate(t, tpl, f)
			f.Close()
			if err != nil {
				return err
			}
		} else if suffix == ".hcl" {
			// Parse team and return TeamVault struct
			tv := parseTeamVault(t.Slug)
			f, err := os.Create(path.Join(dstDirName, fmt.Sprintf("%v-%v%v", tv["ShortName"], tv["Env"], suffix)))
			if err != nil {
				return err
			}
			err = RenderTemplate(tv, tpl, f)
			f.Close()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
