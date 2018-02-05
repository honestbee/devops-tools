package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "rds-snapper"
	app.Usage = "golang tools to manage RDS snapshots"
	app.Version = fmt.Sprintf("1.0.0")

	mainFlag := []cli.Flag{
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
			Name:   "db_name",
			Value:  "",
			Usage:  "origin of snapshots",
			EnvVar: "PLUGIN_DB_NAME",
		},
	}

	exportFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "file",
			Value:  "",
			Usage:  "file to save snapshots list",
			EnvVar: "PLUGIN_FILE",
		},
	}

	clearFlag := []cli.Flag{
		cli.IntFlag{
			Name:   "limit",
			Value:  5,
			Usage:  "number of snapshots to keep",
			EnvVar: "PLUGIN_LIMIT",
		},
	}

	createFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "suffix",
			Usage:  "suffix to add to snapshot name (if not specified, would be a random string)",
			Value:  "",
			EnvVar: "PLUGIN_SUFFIX",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "export",
			Usage: "Export snapshots list to csv file",
			Flags: append(mainFlag, exportFlag...),
			Action: func(c *cli.Context) error {
				file := c.String("file")
				dbName := c.String("db_name")
				accessKey := c.String("aws-access-key")
				secretKey := c.String("aws-secret-key")
				region := c.String("aws-region")
				var err error

				awsConfig := createAwsConfig(accessKey, secretKey, region)
				svc := createRdsClient(awsConfig)
				if dbName != "" {
					err = saveCsv(retrieveInstanceManualSnapshots(dbName, svc), file)
				} else {
					err = saveCsv(retrieveAllManualSnapshots(svc), file)
				}
				return err
			},
		},
		{
			Name:  "clear",
			Usage: "Clear snapshot of specific dbName and only a specified limit number",
			Flags: append(mainFlag, clearFlag...),
			Action: func(c *cli.Context) error {
				limit := c.Int("limit")
				dbName := c.String("db_name")
				accessKey := c.String("aws-access-key")
				secretKey := c.String("aws-secret-key")
				region := c.String("aws-region")

				awsConfig := createAwsConfig(accessKey, secretKey, region)
				svc := createRdsClient(awsConfig)
				maintainSnapshots(dbName, svc, limit)
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Create new snapshot suffix with commit reference",
			Flags: append(mainFlag, createFlag...),
			Action: func(c *cli.Context) error {
				suffix := c.String("suffix")
				dbName := c.String("db_name")
				accessKey := c.String("aws-access-key")
				secretKey := c.String("aws-secret-key")
				region := c.String("aws-region")

				awsConfig := createAwsConfig(accessKey, secretKey, region)
				svc := createRdsClient(awsConfig)
				createSnapshot(dbName, svc, suffix)
				return nil
			},
		},
	}

	return app
}

func main() {

	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
