package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/urfave/cli"
)

// Aws struct type
type Aws struct {
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

// NewAwsClient create *iam client from specific *aws.Config
func NewAwsClient(awsConfig *aws.Config) (*aws.Context, *iam.IAM) {
	ctx := aws.BackgroundContext()
	sess := session.Must(session.NewSession(awsConfig))
	client := iam.New(sess)
	return &ctx, client
}

// ListUsers list all users
func (Aws) ListUsers(c *cli.Context) {
	awsConfig := newAwsConfig(c, c.String("region"))
	ctx, client := NewAwsClient(awsConfig)
	output, err := client.ListUsersWithContext(*ctx, &iam.ListUsersInput{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

// AddUser add new user
func (Aws) AddUser(c *cli.Context) {
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

// RemoveUser delete specific user
func (Aws) RemoveUser(c *cli.Context) {
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

// RemoveUserFromTeams remove user from his / her teams
func (Aws) RemoveUserFromTeams(c *cli.Context) {
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
