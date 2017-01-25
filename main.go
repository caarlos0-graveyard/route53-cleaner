package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/caarlos0/route53-cleaner/internal/addrs"
	"github.com/caarlos0/route53-cleaner/routes"
	"github.com/caarlos0/spin"
	"github.com/urfave/cli"
)

var version = "master"
var debug bool

func isUsed(record *route53.ResourceRecordSet, addrs []string) bool {
	for _, r := range record.ResourceRecords {
		for _, addr := range addrs {
			if *r.Value == addr {
				if debug {
					log.Println(*record.Type, "record to", addr)
				}
				return true
			}
		}
	}
	if record.AliasTarget == nil {
		return false
	}
	var alias = strings.TrimSuffix(*record.AliasTarget.DNSName, ".")
	for _, addr := range addrs {
		if alias == addr {
			if debug {
				log.Println(*record.Type, "ALIAS to", addr)
			}
			return true
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
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "show debug information",
			Destination: &debug,
		},
	}
	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		if debug {
			log.Println("Debug enabled!")
		}
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

		if debug {
			log.Println("Addrs found:", addrs)
			log.Println("\n\n")
		}

		var removables []*route53.ResourceRecordSet
		for _, record := range records {
			if debug {
				log.SetPrefix("")
				log.Println("Checking", *record.Name)
				log.SetPrefix("  --> ")
			}
			if *record.Type != "A" && *record.Type != "AAAA" && *record.Type != "CNAME" {
				if debug {
					log.Println("Type", *record.Type, "ignored")
				}
				continue
			}
			if !isUsed(record, addrs) {
				removables = append(removables, record)
			}
		}
		spin.Stop()

		if debug {
			log.SetPrefix("")
			log.Println("\n")
		}
		for _, record := range removables {
			log.Println(strings.TrimSuffix(*record.Name, "."))
		}

		return nil
	}
	_ = app.Run(os.Args)
}
