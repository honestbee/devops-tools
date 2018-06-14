package main

import (
	"fmt"
	"log"
	"os"

	"github.com/honestbee/devops-tools/noof/pkg/aws"
	"github.com/honestbee/devops-tools/noof/pkg/datadog"
	"github.com/honestbee/devops-tools/noof/pkg/github"
	"github.com/honestbee/devops-tools/noof/pkg/util"
	"github.com/urfave/cli"
)

type Config struct {
	Datadog datadog.Datadog
	Github  github.Github
	Aws     aws.Aws
}

type Action interface {
	AddUser(*cli.Context)
	ListUsers(*cli.Context)
	RemoveUserFromTeams(*cli.Context)
	RemoveUser(*cli.Context)
}

// initApp
func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "common-tools"
	app.Usage = "tool to offload ops tasks for DevOps team"
	app.Version = fmt.Sprintf("0.1.0")

	mainFlag := []cli.Flag{}

	ddFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "datadog-api-key",
			Usage:  "Datadog api key `DATADOG_API_KEY`",
			EnvVar: "PLUGIN_DATADOG_API_KEY,DATADOG_API_KEY",
		},
		cli.StringFlag{
			Name:   "datadog-app-key",
			Usage:  "Datadog app key `DATADOG_APP_KEY`",
			EnvVar: "PLUGIN_DATADOG_APP_KEY,DATADOG_APP_KEY",
		},
		cli.StringFlag{
			Name:  "action",
			Usage: "action",
		},
	}

	githubFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "github-token",
			Usage:  "Github token `GITHUB_TOKEN`",
			EnvVar: "PLUGIN_GITHUB_TOKEN,GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:  "action",
			Usage: "action",
		},
	}

	awsFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key `AWS_ACCESS_KEY`",
			EnvVar: "PLUGIN_ACCESS_KEY,AWS_ACCESS_KEY_ID,AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Key `AWS_SECRET_KEY`",
			EnvVar: "PLUGIN_SECRET_KEY,AWS_SECRET_ACCESS_KEY,AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Value:  "ap-southeast-1",
			Usage:  "AWS Region `AWS_REGION`",
			EnvVar: "PLUGIN_REGION, AWS_REGION",
		},
		cli.StringFlag{
			Name:  "action",
			Usage: "action",
		},
	}

	app.Flags = mainFlag
	app.Commands = []cli.Command{
		{
			Name:   "github",
			Usage:  "manage github users",
			Flags:  githubFlag,
			Action: defaultAction,
		},
		{
			Name:   "datadog",
			Usage:  "manage datadog users",
			Flags:  ddFlag,
			Action: defaultAction,
		},
		{
			Name:   "aws",
			Usage:  "manage aws users",
			Flags:  awsFlag,
			Action: defaultAction,
		},
	}

	return app
}

func defaultAction(c *cli.Context) error {
	action := c.String("action")
	if action == "" {
		log.Fatal("no action provided!")
	}

	// 	if action != "add" && action != "list" && action != "delete" {
	// 		fmt.Println(action)
	// 		log.Fatal("action not valid!")
	// 	}

	var conf Config

	if util.CheckCommand(c.Command.FullName()) == "datadog" {
		executeCommand(conf.Datadog, action, c)
	} else if util.CheckCommand(c.Command.FullName()) == "github" {
		executeCommand(conf.Github, action, c)
	} else {
		executeCommand(conf.Aws, action, c)
	}

	return nil
}

func executeCommand(a Action, action string, c *cli.Context) {
	switch action {
	case "add":
		a.AddUser(c)
	case "list":
		a.ListUsers(c)
	case "RemoveUserFromTeams":
		a.RemoveUserFromTeams(c)
	case "delete":
		a.RemoveUser(c)
	}
}

func main() {

	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
