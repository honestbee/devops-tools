package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func createRdsClient() *rds.RDS {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))

	svc := rds.New(sess)
	return svc
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

func saveCsv(result *rds.DescribeDBSnapshotsOutput) {
	records := [][]string{}

	for index := 0; index < len(result.DBSnapshots); index++ {
		record := result.DBSnapshots[index]
		records = append(records, []string{record.SnapshotCreateTime.String(), *record.DBInstanceIdentifier, *record.DBSnapshotIdentifier})
	}

	outfile, err := os.Create("resultsfile.csv")
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

func main() {

	svc := createRdsClient()

	// result := retrieveSnapshots("hbpay-production", svc)
	// saveCsv(result)
	maintainSnapshots("hbpay-production", svc, 5)
}
