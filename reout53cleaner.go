// Package route53cleaner can find unused entries on Route53
package route53cleaner

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/caarlos0/route53-cleaner/addrs"
	"github.com/caarlos0/route53-cleaner/routes"
)

// UnusedRecord represents a probably unused record that may be removed
type UnusedRecord struct {
	Name, Type, Addr string
}

// FindUnused returns records that are not being used and can probably
// be removed.
func FindUnused(sess *session.Session) (removables []UnusedRecord, err error) {
	records, err := routes.All(sess)
	if err != nil {
		return removables, err
	}
	addrs, err := addrs.All(sess)
	if err != nil {
		return removables, err
	}

	for _, record := range records {
		if *record.Type != "CNAME" {
			continue
		}
		for _, r := range record.ResourceRecords {
			addrs = append(addrs, *r.Value)
		}
	}

	for _, record := range records {
		if *record.Type != "A" && *record.Type != "AAAA" && *record.Type != "CNAME" {
			continue
		}
		if !isUsed(record, addrs) {
			removables = append(removables, newUnused(record))
		}
	}
	return
}

func newUnused(record *route53.ResourceRecordSet) UnusedRecord {
	var addrs []string
	for _, r := range record.ResourceRecords {
		addrs = append(addrs, *r.Value)
	}
	if record.AliasTarget != nil {
		addrs = append(addrs, *record.AliasTarget.DNSName)
	}
	return UnusedRecord{
		Name: *record.Name,
		Addr: strings.Join(addrs, ", "),
		Type: *record.Type,
	}
}

func isUsed(record *route53.ResourceRecordSet, addrs []string) bool {
	for _, r := range record.ResourceRecords {
		for _, addr := range addrs {
			if *r.Value == addr {
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
			return true
		}
	}
	return false
}
