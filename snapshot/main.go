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

func main() {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))

	svc := rds.New(sess)
	input := &rds.DescribeDBSnapshotsInput{
		// DBInstanceIdentifier: aws.String("hbpay-production"),
		SnapshotType:  aws.String("manual"),
		IncludePublic: aws.Bool(true),
		IncludeShared: aws.Bool(true),
	}

	result, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		fmt.Println(err)
	}

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

	// w.Write([]string{"dateCreated", "DBInstanceIdentifier", "DBSnapshotIdentifier"})

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
