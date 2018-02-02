package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func createRdsClient() *rds.RDS {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))

	svc := rds.New(sess)
	return svc
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for index := range b {
		b[index] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createSnapshot(dBInstanceIdentifier string, svc *rds.RDS, suffix string) *rds.CreateDBSnapshotOutput {

	var snapshotName string
	if suffix == "" {
		snapshotName = dBInstanceIdentifier + "-" + randomString(8)
	} else {
		snapshotName = dBInstanceIdentifier + "-" + suffix
	}

	input := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(dBInstanceIdentifier),
		DBSnapshotIdentifier: aws.String(snapshotName),
	}
	result, err := svc.CreateDBSnapshot(input)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func retrieveAllSnapshots(svc *rds.RDS) *rds.DescribeDBSnapshotsOutput {
	input := &rds.DescribeDBSnapshotsInput{
		SnapshotType:  aws.String("manual"),
		IncludePublic: aws.Bool(true),
		IncludeShared: aws.Bool(true),
	}

	result, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func retrieveSnapshots(dBInstanceIdentifier string, svc *rds.RDS) *rds.DescribeDBSnapshotsOutput {
	input := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(dBInstanceIdentifier),
		SnapshotType:         aws.String("manual"),
		IncludePublic:        aws.Bool(true),
		IncludeShared:        aws.Bool(true),
	}

	result, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func cleanUpSnapshots(dBSnapshotIdentifier *string, svc *rds.RDS) {
	input := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(*dBSnapshotIdentifier),
	}

	result, err := svc.DeleteDBSnapshot(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func maintainSnapshots(dBInstanceIdentifier string, svc *rds.RDS, limit int) {
	input := retrieveSnapshots(dBInstanceIdentifier, svc)

	if len(input.DBSnapshots) > limit {
		for index := 0; index < len(input.DBSnapshots)-limit; index++ {
			cleanUpSnapshots(input.DBSnapshots[index].DBSnapshotIdentifier, svc)
		}
	}
}

func saveCsv(result *rds.DescribeDBSnapshotsOutput, filePath string) {
	records := [][]string{}

	for index := 0; index < len(result.DBSnapshots); index++ {
		record := result.DBSnapshots[index]
		records = append(records, []string{record.SnapshotCreateTime.String(), *record.DBInstanceIdentifier, *record.DBSnapshotIdentifier})
	}

	outfile, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Unable to open output")
	}
	defer outfile.Close()

	w := csv.NewWriter(outfile)

	w.Write([]string{"dateCreated", "DBInstanceIdentifier", "DBSnapshotIdentifier"})

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "aws-snapshot-cleanup"
	app.Usage = "golang tools to manage RDS snapshots"
	app.Version = fmt.Sprintf("1.0.%s", build)

	mainFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key `AWS_ACCESS_KEY`",
			EnvVar: "AWS_ACCESS_KEY_ID,AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Key `AWS_SECRET_KEY`",
			EnvVar: "AWS_SECRET_ACCESS_KEY,AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "AWS Region `AWS_REGION`",
			EnvVar: "PLUGIN_AWS_REGION, AWS_REGION",
		},
		cli.StringFlag{
			Name:   "dbName",
			Value:  "",
			Usage:  "origin of snapshots",
			EnvVar: "PLUGIN_DBNAME",
		},
	}

	exportFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "file",
			Usage:  "file to save snapshots list",
			EnvVar: "PLUGIN_FILE",
		},
	}

	clearFlag := []cli.Flag{
		cli.IntFlag{
			Name:   "limit",
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
				dbName := c.String("dbName")

				svc := createRdsClient()
				if dbName != "" {
					saveCsv(retrieveSnapshots(dbName, svc), file)
				} else {
					saveCsv(retrieveAllSnapshots(svc), file)
					fmt.Println("here")
				}
				return nil
			},
		},
		{
			Name:  "clear",
			Usage: "Clear snapshot of specific dbName and only a specified limit number",
			Flags: append(mainFlag, clearFlag...),
			Action: func(c *cli.Context) error {
				limit := c.Int("limit")
				dbName := c.String("dbName")
				svc := createRdsClient()
				maintainSnapshots(dbName, svc, limit)
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Create new snapshot and name it with commit SHA",
			Flags: append(mainFlag, createFlag...),
			Action: func(c *cli.Context) error {
				suffix := c.String("suffix")
				dbName := c.String("dbName")
				svc := createRdsClient()
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
