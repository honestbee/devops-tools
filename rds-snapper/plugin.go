package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type config struct {
	AccessKey string
	SecretKey string
	Region    string
}

// create *aws.Config to use with session
func createAwsConfig(accessKey string, secretKey string, region string) *aws.Config {
	conf := config{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}

	// combine many providers in case some is missing
	creds := credentials.NewChainCredentials([]credentials.Provider{
		// use static access key & private key if available
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     conf.AccessKey,
				SecretAccessKey: conf.SecretKey,
			},
		},
		// fallback to default aws environment variables
		&credentials.EnvProvider{},
		// read aws config file $HOME/.aws/credentials
		&credentials.SharedCredentialsProvider{},
	})

	awsConfig := aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(conf.Region)

	return awsConfig
}

// create *rds.RDS client from specific *aws.Config
func createRdsClient(awsConfig *aws.Config) *rds.RDS {
	sess := session.Must(session.NewSession(awsConfig))
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

func createSnapshot(dbInstanceIdentifier string, svc *rds.RDS, suffix string) *rds.CreateDBSnapshotOutput {

	var snapshotName string
	if suffix == "" {
		snapshotName = dbInstanceIdentifier + "-" + randomString(8)
	} else {
		snapshotName = dbInstanceIdentifier + "-" + suffix
	}

	input := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
		DBSnapshotIdentifier: aws.String(snapshotName),
	}
	result, err := svc.CreateDBSnapshot(input)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func retrieveAllManualSnapshots(svc *rds.RDS) *rds.DescribeDBSnapshotsOutput {
	input := &rds.DescribeDBSnapshotsInput{
		SnapshotType:  aws.String("manual"),
		IncludePublic: aws.Bool(true),
		IncludeShared: aws.Bool(true),
	}

	result, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func retrieveInstanceManualSnapshots(dbInstanceIdentifier string, svc *rds.RDS) *rds.DescribeDBSnapshotsOutput {
	input := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
		SnapshotType:         aws.String("manual"),
		IncludePublic:        aws.Bool(true),
		IncludeShared:        aws.Bool(true),
	}

	result, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func cleanUpSnapshot(dBSnapshotIdentifier *string, svc *rds.RDS) {
	input := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(*dBSnapshotIdentifier),
	}

	result, err := svc.DeleteDBSnapshot(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

// function to maintain specific numbers of snapshot (e.g: limit set to 5, then only keep 5 latest snapshots, delete the others)
func maintainSnapshots(dbInstanceIdentifier string, svc *rds.RDS, limit int) {
	if dbInstanceIdentifier == "" {
		log.Fatal("dbInstanceIdentifier need to be defined!")
	}
	input := retrieveInstanceManualSnapshots(dbInstanceIdentifier, svc)

	if len(input.DBSnapshots) > limit {
		sort.SliceStable(input.DBSnapshots, func(i int, j int) bool {
			return input.DBSnapshots[i].SnapshotCreateTime.Before(*input.DBSnapshots[j].SnapshotCreateTime)
		})
		for index := 0; index < len(input.DBSnapshots)-limit; index++ {
			cleanUpSnapshot(input.DBSnapshots[index].DBSnapshotIdentifier, svc)
		}
	}
}

func saveCsv(result *rds.DescribeDBSnapshotsOutput, filePath string) error {
	records := [][]string{}
	// predefine writer
	var w *csv.Writer

	for index := 0; index < len(result.DBSnapshots); index++ {
		record := result.DBSnapshots[index]
		records = append(records, []string{record.SnapshotCreateTime.String(), *record.DBInstanceIdentifier, *record.DBSnapshotIdentifier})
	}

	if filePath != "" {
		outfile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer outfile.Close()
		w = csv.NewWriter(outfile)
	} else {
		// writing to stdout if filepath is not specified
		w = csv.NewWriter(os.Stdout)
	}

	w.Write([]string{"dateCreated", "DBInstanceIdentifier", "DBSnapshotIdentifier"})

	for _, record := range records {
		if err := w.Write(record); err != nil {
			fmt.Println("error writing record to csv:", err)
			return err
		}
	}
	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if w.Error() != nil {
		return w.Error()
	}
	return nil
}
