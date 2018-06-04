package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	iam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/go-github/github"
	"github.com/honestbee/devops-tools/noof/pkg/util"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	dd "gopkg.in/zorkian/go-datadog-api.v2"
)

type Config struct {
	Datadog Datadog
	Github  Github
	Aws     Aws
}

type Datadog struct {
}

type Github struct {
}

type Aws struct {
}

type Action interface {
	addUser(*cli.Context)
	listUsers(*cli.Context)
	removeUserFromTeams(*cli.Context)
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

func NewDatadogClient(c *cli.Context) *dd.Client {
	return dd.NewClient(c.String("datadog-api-key"), c.String("datadog-app-key"))
}

func NewGithubClient(c *cli.Context) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("github-token")})
	tc := oauth2.NewClient(ctx, ts)
	return ctx, github.NewClient(tc)
}

// create *aws.Config to use with session
func newAwsConfig(c *cli.Context, region string) *aws.Config {
	// combine many providers in case some is missing
	creds := credentials.NewChainCredentials([]credentials.Provider{
		// use static access key & private key if available
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     c.String("aws-access-key"),
				SecretAccessKey: c.String("aws-secret-key"),
			},
		},
		// fallback to default aws environment variables
		&credentials.EnvProvider{},
		// read aws config file $HOME/.aws/credentials
		&credentials.SharedCredentialsProvider{},
	})

	awsConfig := aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(region)

	return awsConfig
}

// create *iam client from specific *aws.Config
func NewAwsClient(awsConfig *aws.Config) (*aws.Context, *iam.IAM) {
	ctx := aws.BackgroundContext()
	sess := session.Must(session.NewSession(awsConfig))
	client := iam.New(sess)
	return &ctx, client
}

func executeCommand(a Action, action string, c *cli.Context) {
	switch action {
	case "add":
		a.addUser(c)
	case "list":
		a.listUsers(c)
	case "removeUserFromTeams":
		a.removeUserFromTeams(c)
	case "delete":
		a.deleteUser(c)
	}
}

func (d Datadog) listUsers(c *cli.Context) {
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

func (d Datadog) listUserTeams(c *cli.Context) {

}

func (d Datadog) removeUserFromTeams(c *cli.Context) {

}

func (g Github) addUser(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	_, _, err := client.Organizations.EditOrgMembership(ctx, c.Args().Get(0), "honestbee", &github.Membership{})
	if err != nil {
		fmt.Println(err)
	}
}

func (g Github) listUsers(c *cli.Context) {
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

func listUserTeams(c *cli.Context, ctx context.Context, client *github.Client) []*github.Team {
	// https://godoc.org/github.com/google/go-github/github#OrganizationsService.ListUserTeams
	var teamList []*github.Team
	teams, _, err := client.Organizations.ListTeams(ctx, "honestbee", &github.ListOptions{})
	for _, team := range teams {
		isTeamMember, _, err := client.Organizations.IsTeamMember(ctx, team.GetID(), c.Args().Get(0))
		if err != nil {
			fmt.Println(err)
		}
		if isTeamMember {
			fmt.Println(*team.Name)
			teamList = append(teamList, team)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return teamList
}

func (g Github) removeUserFromTeams(c *cli.Context) {
	// https://godoc.org/github.com/google/go-github/github#OrganizationsService.RemoveTeamMembership
	ctx, client := NewGithubClient(c)
	teams := listUserTeams(c, ctx, client)
	for _, team := range teams {
		_, err := client.Organizations.RemoveTeamMembership(ctx, team.GetID(), c.Args().Get(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (Aws) listUsers(c *cli.Context) {
	awsConfig := newAwsConfig(c, c.String("region"))
	ctx, client := NewAwsClient(awsConfig)
	output, err := client.ListUsersWithContext(*ctx, &iam.ListUsersInput{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func (Aws) addUser(c *cli.Context) {
	awsConfig := newAwsConfig(c, c.String("region"))
	ctx, client := NewAwsClient(awsConfig)
	userName := c.Args().Get(0)
	output, err := client.CreateUserWithContext(*ctx, &iam.CreateUserInput{
		UserName: &userName,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output.User.UserName)
}

func (Aws) deleteUser(c *cli.Context) {
	awsConfig := newAwsConfig(c, c.String("region"))
	ctx, client := NewAwsClient(awsConfig)
	userName := c.Args().Get(0)
	output, err := client.DeleteUserWithContext(*ctx, &iam.DeleteUserInput{
		UserName: &userName,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func awsListUserTeams(c *cli.Context, awsConfig *aws.Config, ctx *aws.Context, client *iam.IAM) []*iam.Group {
	var groupList []*iam.Group
	userName := c.Args().Get(0)
	output, err := client.ListGroupsForUserWithContext(*ctx, &iam.ListGroupsForUserInput{
		UserName: &userName,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, group := range output.Groups {
		fmt.Println(*group.GroupName)
		groupList = append(groupList, group)
	}

	return groupList
}
func (Aws) removeUserFromTeams(c *cli.Context) {
	awsConfig := newAwsConfig(c, c.String("region"))
	ctx, client := NewAwsClient(awsConfig)
	userName := c.Args().Get(0)
	groups := awsListUserTeams(c, awsConfig, ctx, client)

	for _, group := range groups {
		client.RemoveUserFromGroupWithContext(*ctx, &iam.RemoveUserFromGroupInput{
			UserName:  &userName,
			GroupName: group.GroupName,
		})
	}
}

func main() {

	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
