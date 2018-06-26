package drone

import (
	"fmt"
	"os"
	"strings"

	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

// NewClient returns a new client from the CLI context.
func NewClient() (drone.Client, error) {
	var (
		token  = os.Getenv("DRONE_TOKEN")
		server = os.Getenv("DRONE_SERVER")
	)
	server = strings.TrimRight(server, "/")

	// if no server url is provided we can default
	// to the hosted Drone service.
	if len(server) == 0 {
		return nil, fmt.Errorf("Error: you must provide the Drone server address")
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("Error: you must provide your Drone access token")
	}

	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)

	return drone.NewClient(server, auther), nil
}
