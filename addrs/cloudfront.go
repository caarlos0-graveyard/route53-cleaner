package addrs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func cloudFrontResolver(sess *session.Session, cfg *aws.Config) (addrs []string, err error) {
	err = cloudfront.New(sess, cfg).ListDistributionsPages(
		&cloudfront.ListDistributionsInput{},
		func(output *cloudfront.ListDistributionsOutput, lastPage bool) bool {
			for _, dist := range output.DistributionList.Items {
				if dist.DomainName != nil {
					addrs = append(addrs, *dist.DomainName)
				}
				if dist.Aliases == nil {
					continue
				}
				for _, alias := range dist.Aliases.Items {
					addrs = append(addrs, *alias)
				}
			}
			return !lastPage
		},
	)
	return
}
