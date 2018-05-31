package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	dd "gopkg.in/zorkian/go-datadog-api.v2"
)

type Config struct {
	Datadog Datadog
}

type Datadog struct {
	APIKey string
	AppKey string
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
	}

	app.Flags = mainFlag

	app.Commands = []cli.Command{
		{
			Name:  "datadog",
			Usage: "manage datadog users",
			Flags: ddFlag,
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add new user",
					Action: func(c *cli.Context) error {
						addUser(c)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list all users",
					Action: func(c *cli.Context) error {
						listUser(c)
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "delete user",
					Action: func(c *cli.Context) error {
						deleteUser(c)
						return nil
					},
				},
			},
		},
	}

	return app
}

func envAction(c *cli.Context) Config {
	var conf Config
	conf.Datadog.APIKey = os.Getenv("DATADOG_API_KEY")
	conf.Datadog.AppKey = os.Getenv("DATADOG_APP_KEY")

	return conf
}

func clientAction(c *cli.Context, conf Config) *dd.Client {
	client := dd.NewClient(conf.Datadog.APIKey, conf.Datadog.AppKey)

	return client
}

func listUser(c *cli.Context) {
	client := clientAction(c, envAction(c))
	users, _ := client.GetUsers()

	for _, user := range users {
		fmt.Println(*user.Email)
	}
}

func addUser(c *cli.Context) {
	username := c.Args().First()
	client := clientAction(c, envAction(c))
	user, _ := client.CreateUser(&username, &username)

	fmt.Println(*user.Email)

}

func deleteUser(c *cli.Context) {
	username := c.Args().First()
	client := clientAction(c, envAction(c))
	client.DeleteUser(username)

}

func main() {

	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
