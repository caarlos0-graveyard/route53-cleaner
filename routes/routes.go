// Package routes can find all route53 records of a given account
package routes

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// All route53 records of the account
func All(sess *session.Session) (records []*route53.ResourceRecordSet, err error) {
	svc := route53.New(sess)

	var zones []*route53.HostedZone
	if err := svc.ListHostedZonesPages(
		&route53.ListHostedZonesInput{},
		func(output *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
			zones = append(zones, output.HostedZones...)
			return !lastPage
		},
	); err != nil {
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
