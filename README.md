# Route53 Cleaner

Suggests records that could be deleted from your AWS Route53 hosted zones.

[![Release](https://img.shields.io/github/release/caarlos0/route53-cleaner.svg?style=flat-square)](https://github.com/caarlos0/route53-cleaner/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/caarlos0/route53-cleaner.svg?style=flat-square)](https://travis-ci.org/caarlos0/route53-cleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/route53-cleaner?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/route53-cleaner)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

## Install

```console
$ brew install caarlos0/tap/route53-cleaner
```

## How it works

Route53 Cleaner scans all your instances, ELBs, RDSs and Route53 zones' records,
then compiles and prints a list of records that might not being used.

Please note that Route53 Cleaner **will never change anything in your account**. You
can check the code or give it read-only keys if you do not trust that affirmation 
(I won't blame you for that).

## Auth

Either by having a `~/.aws/credentials`, `~/.aws/config` or the `AWS_ACCESS_KEY_ID` and 
`AWS_SECRET_ACCESS_KEY` environment variables exported.

More info [here](
)

## Contributing

Please refer to our [contributing guidelines](CONTRIBUTING.md).
