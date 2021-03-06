package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	route53cleaner "github.com/caarlos0/route53-cleaner"
	"github.com/caarlos0/spin"
	"github.com/urfave/cli"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Name = "route53-cleaner"
	app.Version = version
	app.Author = "Carlos Alexandro Becker (root@carlosbecker.com)"
	app.Usage = "Find records that could be deleted from your AWS Route53 hosted zones"
	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		var spin = spin.New("\033[36m %s Working...\033[m")
		spin.Start()
		defer spin.Stop()
		sess, err := session.NewSession()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		removables, err := route53cleaner.FindUnused(sess)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		spin.Stop()
		for _, record := range removables {
			log.Println(strings.TrimSuffix(record.Name, "."), record.Type, record.Addr)
		}

		return nil
	}
	_ = app.Run(os.Args)
}
