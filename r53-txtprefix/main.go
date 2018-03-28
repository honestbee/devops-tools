package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "hosted-zone-id",
			Usage:  "Hosted zone `id` for route53",
			EnvVar: "AWS_HOSTED_ZONE_ID",
		},
		cli.StringFlag{
			Name:  "prefix",
			Value: "prefix.",
			Usage: "`prefix` for external-dns TXT records",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "error",
			Usage:  "Log level (panic, fatal, error, warn, info, or debug)",
			EnvVar: "LOG_LEVEL",
		},
		cli.StringFlag{
			Name:   "region",
			Value:  "ap-southeast-1",
			Usage:  "default `region` to use",
			EnvVar: "AWS_REGION",
		},
		cli.StringFlag{
			Name:  "access-key-id",
			Value: "",
			Usage: "default `key` to use",
		},
		cli.StringFlag{
			Name:  "secret-access-key",
			Value: "",
			Usage: "default `secret` to use",
		},
	}
	app := cli.NewApp()
	app.Name = "r53-txtprefix"
	app.Usage = "Prefix all r53 of type TXT with provided string"
	app.Action = run

	app.Version = fmt.Sprintf("0.1.%s", build)
	app.Author = "Honestbee DevOps"

	app.Flags = flags

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func run(c *cli.Context) error {
	logLevelString := c.String("log-level")
	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	awsRegion := c.String("region")
	awsAccessKey := c.String("access-key-id")
	awsSecretKey := c.String("secret-access-key")
	awsHostedZoneID := aws.String(c.String("hosted-zone-id"))
	if *awsHostedZoneID == "" {
		return fmt.Errorf("hosted-zone-id is mandatory")
	}

	prefix := c.String("prefix")

	conf := initAwsConfig(awsRegion, awsAccessKey, awsSecretKey)
	client := route53.New(session.New(conf))
	ctx := context.Background()
	err = client.ListResourceRecordSetsPagesWithContext(ctx,
		&route53.ListResourceRecordSetsInput{
			HostedZoneId: awsHostedZoneID,
		},
		func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, rrset := range page.ResourceRecordSets {
				if *rrset.Type == route53.RRTypeTxt {
					if fromExternalDNS(rrset) {
						if !strings.HasPrefix(*rrset.Name, prefix) {
							log.Infof("Prefixing %q ...", *rrset.Name)
							newName := aws.String(prefix + *rrset.Name)
							changeRRSet := &route53.ChangeResourceRecordSetsInput{
								ChangeBatch: &route53.ChangeBatch{
									Changes: []*route53.Change{
										{
											Action:            aws.String("DELETE"),
											ResourceRecordSet: rrset,
										},
										{
											Action: aws.String("CREATE"),
											ResourceRecordSet: &route53.ResourceRecordSet{
												Name:            newName,
												ResourceRecords: rrset.ResourceRecords,
												TTL:             rrset.TTL,
												Type:            rrset.Type,
											},
										},
									},
									Comment: aws.String("Update prefix"),
								},
								HostedZoneId: awsHostedZoneID,
							}
							log.Debugf("Change %v", changeRRSet)
							result, err := client.ChangeResourceRecordSets(changeRRSet)
							if err != nil {
								log.Errorf("Error prefixing %q: %s", *rrset.Name, err)
							}
							log.Debug(result)
							log.Infof("Prefixed %q to %q!", *rrset.Name, *newName)
						} else {
							log.Infof("Skipping (already prefixed) %q", *rrset.Name)
						}
					}
				}
			}
			return true // page until end
		})

	return err
}

func fromExternalDNS(rrset *route53.ResourceRecordSet) bool {
	for _, rr := range rrset.ResourceRecords {
		if strings.Contains(*rr.Value, "heritage=external-dns") {
			return true
		}
	}
	return false
}

func initAwsConfig(region, accessKey, secretKey string) *aws.Config {
	awsConfig := aws.NewConfig()
	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(region)
	return awsConfig
}
