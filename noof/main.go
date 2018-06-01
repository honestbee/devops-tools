package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/honestbee/devops-tools/noof/pkg/util"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	dd "gopkg.in/zorkian/go-datadog-api.v2"
)

type Config struct {
	Datadog Datadog
	Github  Github
}

type Datadog struct {
}

type Github struct {
}

type Action interface {
	addUser(*cli.Context)
	listUser(*cli.Context)
	deleteUser(*cli.Context)
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
	}

	return app
}

func NewDatadogClient(c *cli.Context) *dd.Client {
	return dd.NewClient(c.String("datadog-api-key"), c.String("datadog-app-key"))

}

func NewGithubClient(c *cli.Context) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("github-token")})
	tc := oauth2.NewClient(ctx, ts)
	return ctx, github.NewClient(tc)
}

func defaultAction(c *cli.Context) error {
	action := c.String("action")
	if action == "" {
		log.Fatal("no action provided!")
	}

	if action != "add" && action != "list" && action != "delete" {
		fmt.Println(action)
		log.Fatal("action not valid!")
	}

	var conf Config

	if util.CheckCommand(c.Command.FullName()) == "datadog" {
		executeCommand(conf.Datadog, action, c)
	} else {
		executeCommand(conf.Github, action, c)
	}

	return nil
}

func executeCommand(a Action, action string, c *cli.Context) {
	switch action {
	case "add":
		a.addUser(c)
	case "list":
		a.listUser(c)
	case "delete":
		a.deleteUser(c)
	}
}

func (d Datadog) listUser(c *cli.Context) {
	client := NewDatadogClient(c)
	users, _ := client.GetUsers()

	for _, user := range users {
		fmt.Println(*user.Email)
	}
}

func (d Datadog) addUser(c *cli.Context) {
	username := c.Args().Get(0)
	client := NewDatadogClient(c)
	user, _ := client.CreateUser(&username, &username)

	fmt.Println(*user.Email)

}

func (d Datadog) deleteUser(c *cli.Context) {
	username := c.Args().Get(0)
	client := NewDatadogClient(c)
	client.DeleteUser(username)
}

func (g Github) addUser(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	_, _, err := client.Organizations.EditOrgMembership(ctx, c.Args().Get(0), "honestbee", &github.Membership{})
	if err != nil {
		fmt.Println(err)
	}
}

func (g Github) listUser(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	users, _, err := client.Organizations.ListMembers(ctx, "honestbee", &github.ListMembersOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1000,
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range users {
		fmt.Println(*user.Login)
	}
}

func (g Github) deleteUser(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	_, err := client.Organizations.RemoveOrgMembership(ctx, c.Args().Get(0), "honestbee")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
