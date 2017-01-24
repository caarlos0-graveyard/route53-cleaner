package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/caarlos0/spin"
	"github.com/urfave/cli"
)

var version = "master"

func allRecordSets(sess *session.Session) (records []*route53.ResourceRecordSet, err error) {
	svc := route53.New(sess)

	var zones []*route53.HostedZone
	if err := svc.ListHostedZonesPages(
		&route53.ListHostedZonesInput{},
		func(output *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
			zones = append(zones, output.HostedZones...)
			return *output.IsTruncated
		},
	);
		err != nil {
		return records, err
	}

	for _, zone := range zones {
		if err := svc.ListResourceRecordSetsPages(&route53.ListResourceRecordSetsInput{
			HostedZoneId: zone.Id,
		}, func(output *route53.ListResourceRecordSetsOutput, lastPage bool) (shouldContinue bool) {
			records = append(records, output.ResourceRecordSets...)
			return *output.IsTruncated
		}); err != nil {
			return records, err
		}
	}
	return
}

func allValidAddrs(sess *session.Session) (addrs []string, err error) {
	return
}

func main() {
	app := cli.NewApp()
	app.Name = "route53-cleaner"
	app.Version = version
	app.Author = "Carlos Alexandro Becker (caarlos0@gmail.com)"
	app.Usage = "Find possibly unused route53 records"
	app.Flags = []cli.Flag{}
	app.Action = func(c *cli.Context) error {
		sess, err := session.NewSession()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		spin := spin.New("\033[36m %s Working...\033[m")
		spin.Start()
		records, err := allRecordSets(sess)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		//var ips []string

		spin.Stop()

		var removables []*route53.ResourceRecordSet
		for _, record := range records {
			if len(record.ResourceRecords) == 0 && record.AliasTarget == nil {
				removables = append(removables, record)
			}
			//fmt.Println(*record.Name, "can be removed")
			fmt.Println(record)
		}
		return nil
	}
	app.Run(os.Args)
}
