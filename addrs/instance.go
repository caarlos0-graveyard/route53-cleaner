package addrs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func instanceResolver(sess *session.Session, cfg *aws.Config) (addrs []string, err error) {
	err = ec2.New(sess, cfg).DescribeInstancesPages(
		&ec2.DescribeInstancesInput{},
		func(output *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range output.Reservations {
				for _, instance := range reservation.Instances {
					// TODO: maybe avoid terminated-ing instances?
					for _, eni := range instance.NetworkInterfaces {
						addrs = append(addrs, *eni.PrivateIpAddress)
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
