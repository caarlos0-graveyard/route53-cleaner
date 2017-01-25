package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/caarlos0/route53-cleaner/internal/addrs"
	"github.com/caarlos0/route53-cleaner/routes"
	"github.com/caarlos0/spin"
	"github.com/urfave/cli"
)

var version = "master"

func isUsed(record *route53.ResourceRecordSet, addrs []string) bool {
	for _, r := range record.ResourceRecords {
		for _, addr := range addrs {
			if *r.Value == addr {
				return true
			}
		}
	}
	return false
}

func main() {
	app := cli.NewApp()
	app.Name = "route53-cleaner"
	app.Version = version
	app.Author = "Carlos Alexandro Becker (root@carlosbecker.com)"
	app.Usage = "Find records that could be deleted from your AWS Route53 hosted zones"
	app.Action = func(c *cli.Context) error {
		sess, err := session.NewSession()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		spin := spin.New("\033[36m %s Working...\033[m")
		spin.Start()

		records, err := routes.All(sess)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		addrs, err := addrs.All(sess)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, record := range records {
			if *record.Type != "CNAME" {
				continue
			}
			for _, r := range record.ResourceRecords {
				addrs = append(addrs, *r.Value)
			}
		}

		var removables []*route53.ResourceRecordSet
		for _, record := range records {
			if *record.Type != "A" && *record.Type != "AAAA" && *record.Type != "CNAME" {
				continue
			}
			if !isUsed(record, addrs) {
				removables = append(removables, record)
			}
		}
		spin.Stop()

		for _, record := range removables {
			fmt.Println(*record.Name, "might be removed")
		}

		return nil
	}
	_ = app.Run(os.Args)
}
