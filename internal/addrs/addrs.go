package addrs

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/sync/errgroup"
)

type resolver func(sess *session.Session, cfg *aws.Config) (addrs []string, err error)

func All(sess *session.Session) (addrs []string, err error) {
	resp, err := ec2.New(sess, aws.NewConfig().WithRegion("us-east-1")).DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return
	}
	var resolvers = []resolver{
		instanceResolver,
		rdsResolver,
		elbResolver,
	}
	var g errgroup.Group
	var m sync.Mutex
	for _, region := range resp.Regions {
		var cfg = aws.NewConfig().WithRegion(*region.RegionName)
		for _, resolver := range resolvers {
			resolver := resolver
			g.Go(func() (err error) {
				result, err := resolver(sess, cfg)
				if err == nil {
					m.Lock()
					defer m.Unlock()
					addrs = append(addrs, result...)
				}
				return
			})
		}
	}
	err = g.Wait()
	return
}
