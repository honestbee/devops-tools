package datadog

import (
	"fmt"

	"github.com/urfave/cli"
	dd "gopkg.in/zorkian/go-datadog-api.v2"
)

type Datadog struct {
}

func NewDatadogClient(c *cli.Context) *dd.Client {
	return dd.NewClient(c.String("datadog-api-key"), c.String("datadog-app-key"))
}

func (d Datadog) ListUsers(c *cli.Context) {
	client := NewDatadogClient(c)
	users, _ := client.GetUsers()

	for _, user := range users {
		fmt.Println(*user.Email)
	}
}

func (d Datadog) AddUser(c *cli.Context) {
	username := c.Args().Get(0)
	client := NewDatadogClient(c)
	user, _ := client.CreateUser(&username, &username)

	fmt.Println(*user.Email)

}

func (d Datadog) DeleteUser(c *cli.Context) {
	username := c.Args().Get(0)
	client := NewDatadogClient(c)
	client.DeleteUser(username)
}

func (d Datadog) listUserTeams(c *cli.Context) {

}

func (d Datadog) RemoveUserFromTeams(c *cli.Context) {

}
