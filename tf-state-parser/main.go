package main

import (
	"fmt"
	"os"

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
		cli.StringFlag{
			Name:  "source, s",
			Value: "state.json",
			Usage: "Source `state` file",
		},
		// cli.StringFlag{
		// 	Name:  "destination, d",
		// 	Value: "output/teams-config",
		// 	Usage: "Destination `directory` to render in - must exist",
		// },
		// cli.StringFlag{
		// 	Name:  "template, t",
		// 	Value: "templates/team.tf.tpl",
		// 	Usage: "Desired template used to render output",
		// },
	}
	app := cli.NewApp()
	app.Name = "tf-state-parser"
	app.Usage = "Parse TF state"
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

	var s state
	stateFile, err := os.Open(c.String("source"))
	defer stateFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening tfstate file: %s\n", err)
		return err
	}

	err = s.read(stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading tfstate file: %s\n", err)
		return err
	}

	if s.Modules == nil {
		fmt.Fprint(os.Stderr, "No modules in tfstate file")
		return fmt.Errorf("No modules in tfstate file")
	}

	writeImports(&s, &gitHub)

	return nil
}
