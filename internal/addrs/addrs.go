package addrs

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/rds"
	"golang.org/x/sync/errgroup"
)

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
			g.Go(func() error {
				result, err := resolver(sess, cfg)
				if err != nil {
					return err
				}
				m.Lock()
				defer m.Unlock()
				addrs = append(addrs, result...)
				return nil
			})
		}
	}
	err = g.Wait()
	return
}

type resolver func(sess *session.Session, cfg *aws.Config) (addrs []string, err error)

func instanceResolver(sess *session.Session, cfg *aws.Config) (addrs []string, err error) {
	err = ec2.New(sess, cfg).DescribeInstancesPages(
		&ec2.DescribeInstancesInput{},
		func(output *ec2.DescribeInstancesOutput, lastPage bool) bool {
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
	)
	return
}

func elbResolver(sess *session.Session, cfg *aws.Config) (addrs []string, err error) {
	err = elb.New(sess, cfg).DescribeLoadBalancersPages(
		&elb.DescribeLoadBalancersInput{},
		func(output *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, e := range output.LoadBalancerDescriptions {
				if e.DNSName != nil {
					addrs = append(addrs, *e.DNSName)
				}
			}
			return !lastPage
		},
	)
	return
}

func rdsResolver(sess *session.Session, cfg *aws.Config) (addrs []string, err error) {
	err = rds.New(sess, cfg).DescribeDBInstancesPages(
		&rds.DescribeDBInstancesInput{},
		func(output *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			for _, db := range output.DBInstances {
				addrs = append(addrs, *db.Endpoint.Address)
			}
			return !lastPage
		},
	)
	return
}
