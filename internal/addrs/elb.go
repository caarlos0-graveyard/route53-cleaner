package addrs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

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
