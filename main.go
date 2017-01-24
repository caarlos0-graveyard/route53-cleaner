package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/caarlos0/spin"
	"github.com/urfave/cli"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

var version = "master"

func allRecordSets(sess *session.Session) (records []*route53.ResourceRecordSet, err error) {
	svc := route53.New(sess)

	var zones []*route53.HostedZone
	if err := svc.ListHostedZonesPages(
		&route53.ListHostedZonesInput{},
		func(output *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
			zones = append(zones, output.HostedZones...)
			return !lastPage
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
			return !lastPage
		}); err != nil {
			return records, err
		}
	}
	return
}

func allValidAddrs(sess *session.Session) (addrs []string, err error) {
	resp, err := ec2.New(sess, aws.NewConfig().WithRegion("us-east-1")).DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return
	}
	for _, region := range resp.Regions {
		var cfg = aws.NewConfig().WithRegion(*region.RegionName)
		if err := ec2.New(sess, cfg).DescribeInstancesPages(
			&ec2.DescribeInstancesInput{},
			func(output *ec2.DescribeInstancesOutput, lastPage bool) (shouldContinue bool) {
				for _, reservation := range output.Reservations {
					for _, instance := range reservation.Instances {
						// TODO: maybe avoid terminated-ing instances?
						if instance.PrivateIpAddress != nil {
							addrs = append(addrs, *instance.PrivateIpAddress)
						}
						if instance.PublicDnsName != nil {
							addrs = append(addrs, *instance.PublicDnsName)
						}
						if instance.PublicIpAddress != nil {
							addrs = append(addrs, *instance.PublicIpAddress)
						}
					}
				}
				return !lastPage
			},
		); err != nil {
			return addrs, err
		}
		if err := elbv2.New(sess, cfg).DescribeLoadBalancersPages(
			&elbv2.DescribeLoadBalancersInput{},
			func(output *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
				for _, elb := range output.LoadBalancers {
					if elb.DNSName != nil {
						addrs = append(addrs, *elb.DNSName)
					}
				}
				return !lastPage
			},
		); err != nil {
			return addrs, err
		}
	}
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

		addrs, err := allValidAddrs(sess)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		records, err := allRecordSets(sess)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		var removables []*route53.ResourceRecordSet
		for _, record := range records {
			if *record.Type == "NS" || *record.Type == "SOA" {
				continue
			}
			var used bool
			for _, r := range record.ResourceRecords {
				for _, addr := range addrs {
					if *r.Value == addr {
						used = true
						continue
					}
				}
			}
			if !used {
				removables = append(removables, record)
			}
		}
		spin.Stop()

		for _, record := range removables {
			fmt.Println(*record.Name, "might be removed")
		}

		return nil
	}
	app.Run(os.Args)
}
