package main

import (
	"fmt"
	"log"
	"os"

	"github.com/honestbee/devops-tools/noof/pkg/util"
	"github.com/urfave/cli"
	dd "gopkg.in/zorkian/go-datadog-api.v2"
)

type Config struct {
	Datadog Datadog
	Github  Github
}

type Datadog struct {
	APIKey string
	AppKey string
}

type Github struct {
	APIKey string
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
			Name:   "github-api-key",
			Usage:  "Github api key `GITHUB_API_KEY`",
			EnvVar: "PLUGIN_GITHUB_API_KEY,GITHUB_API_KEY",
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
	var conf Config
	conf.Datadog.APIKey = os.Getenv("DATADOG_API_KEY")
	conf.Datadog.AppKey = os.Getenv("DATADOG_APP_KEY")
	client := dd.NewClient(conf.Datadog.APIKey, conf.Datadog.AppKey)

	return client
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

	if util.CheckCommand(c.Command.FullName()) == "datadog" {
		var d Datadog
		executeCommand(d, action, c)
	} else {
		var g Github
		executeCommand(g, action, c)
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
	username := c.Args().First()
	client := NewDatadogClient(c)
	user, _ := client.CreateUser(&username, &username)

	fmt.Println(*user.Email)

}

func (d Datadog) deleteUser(c *cli.Context) {
	username := c.Args().First()
	client := NewDatadogClient(c)
	client.DeleteUser(username)

}

func (g Github) addUser(c *cli.Context) {
	fmt.Println("hello world")
}

func (g Github) listUser(c *cli.Context) {
	fmt.Println("hello world")
}

func (g Github) deleteUser(c *cli.Context) {
	fmt.Println("hello world")
}

func main() {

	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
