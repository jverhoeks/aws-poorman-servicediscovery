package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var record string
var value string
var recordtype string
var TTL int64
var weight = int64(1)
var zoneId string

func init() {
	flag.StringVar(&record, "r", "", "dns record")
	flag.StringVar(&value, "v", "", "value of dns record")
	flag.StringVar(&recordtype, "t", "A", "dns record type")
	flag.StringVar(&zoneId, "z", "", "AWS Zone Id for domain")
	flag.Int64Var(&TTL, "ttl", int64(60), "TTL for DNS Cache")

}

func main() {
	flag.Parse()

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	if value == "" {
		getLocalIP(sess)
	}

	if record == "" || zoneId == "" || value == "" {
		fmt.Println(fmt.Errorf("Incomplete arguments: r: %s, v: %s, z: %s\n", record, value, zoneId))
		flag.PrintDefaults()
		return
	}

	svc := route53.New(sess)
	updateRecord(svc)
}

func getLocalIP(sess *session.Session) error {
	// Create a EC2Metadata client from just a session.
	svc := ec2metadata.New(sess)
	if svc.Available() {
		localipv4, err := svc.GetMetadata("local-ipv4")
		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println("Error: " + err.Error())
			return err
		}
		value = localipv4
		fmt.Println("Retrieved local ip from aws metadata : " + value)
	} else {
		fmt.Println("Error: no metadata available (not running on ecs?), provide value as parameter")
		return nil
	}
	return nil
}

func updateRecord(svc *route53.Route53) {

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(record),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(value),
							},
						},
						TTL:  aws.Int64(TTL),
						Type: aws.String(recordtype),
					},
				},
			},
			Comment: aws.String("Updating from update-route53.go"),
		},
		HostedZoneId: aws.String(zoneId), // Required
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println("Change Response:")
	fmt.Println(resp)
}
