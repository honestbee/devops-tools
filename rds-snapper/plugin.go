package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

// config struct to use with AWS
type config struct {
	AccessKey string
	SecretKey string
	Region    string
	Limit     int
	DbName    string
	Suffix    string
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// create *aws.Config to use with session
func createAwsConfig(accessKey string, secretKey string, region string) *aws.Config {
	// combine many providers in case some is missing
	creds := credentials.NewChainCredentials([]credentials.Provider{
		// use static access key & private key if available
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
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

// create *rds.RDS client from specific *aws.Config
func createRdsClient(awsConfig *aws.Config) *rds.RDS {
	sess := session.Must(session.NewSession(awsConfig))
	svc := rds.New(sess)
	return svc
}

// https://stackoverflow.com/a/31832326/2490986
func randomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
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

	fmt.Print(snapshotName, " has been created successfully!")
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

func deleteSnapshot(dBSnapshotIdentifier *string, svc *rds.RDS) {
	input := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(*dBSnapshotIdentifier),
	}

	result, err := svc.DeleteDBSnapshot(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

// function to maintain specific numbers of snapshot (e.g: keep set to 5, then only keep 5 latest snapshots, delete the others)
func clearSnapshots(dbInstanceIdentifier string, svc *rds.RDS, keep int) {
	if dbInstanceIdentifier == "" {
		log.Fatal("dbInstanceIdentifier need to be defined!")
	}
	input := retrieveInstanceManualSnapshots(dbInstanceIdentifier, svc)

	if len(input.DBSnapshots) > keep && keep > 0 {
		sort.SliceStable(input.DBSnapshots, func(i int, j int) bool {
			return input.DBSnapshots[i].SnapshotCreateTime.Before(*input.DBSnapshots[j].SnapshotCreateTime)
		})
		for index := 0; index < len(input.DBSnapshots)-keep; index++ {
			deleteSnapshot(input.DBSnapshots[index].DBSnapshotIdentifier, svc)
		}
	}
}

func createWriter(output string) (*csv.Writer, error) {
	if output != "" {
		outfile, err := os.Create(output)
		if err != nil {
			return nil, err
		}
		defer outfile.Close()
		return csv.NewWriter(outfile), nil
	}
	return csv.NewWriter(os.Stdout), nil
}

func saveCsv(result *rds.DescribeDBSnapshotsOutput, w *csv.Writer) error {
	records := [][]string{}

	for index := 0; index < len(result.DBSnapshots); index++ {
		record := result.DBSnapshots[index]
		records = append(records, []string{record.SnapshotCreateTime.String(), *record.DBInstanceIdentifier, *record.DBSnapshotIdentifier})
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
