package datadog

import (
	"fmt"

	"github.com/urfave/cli"
	dd "gopkg.in/zorkian/go-datadog-api.v2"
)

// Datadog data type
type Datadog struct {
}

// NewDatadogClient create new datadog client
func NewDatadogClient(c *cli.Context) *dd.Client {
	return dd.NewClient(c.String("datadog-api-key"), c.String("datadog-app-key"))
}

// ListUsers list all users
func (d Datadog) ListUsers(c *cli.Context) {
	client := NewDatadogClient(c)
	users, _ := client.GetUsers()

	for _, user := range users {
		fmt.Println(*user.Email)
	}
}

// AddUser add new user
func (d Datadog) AddUser(c *cli.Context) {
	username := c.Args().Get(0)
	client := NewDatadogClient(c)
	client.CreateUser(&username, &username)

}

// RemoveUser remove user from datadog
func (d Datadog) RemoveUser(c *cli.Context) {
	username := c.Args().Get(0)
	client := NewDatadogClient(c)
	client.DeleteUser(username)
}

func (d Datadog) listUserTeams(c *cli.Context) {

}

// RemoveUserFromTeams remove user from his teams
func (d Datadog) RemoveUserFromTeams(c *cli.Context) {

}
