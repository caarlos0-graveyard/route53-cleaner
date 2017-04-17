package addrs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

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
