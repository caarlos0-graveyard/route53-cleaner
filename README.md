# Route53 Cleaner

Suggests records that could be deleted from your AWS Route53 hosted zones.

[![Release](https://img.shields.io/github/release/caarlos0/route53-cleaner.svg?style=flat-square)](https://github.com/caarlos0/route53-cleaner/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/caarlos0/route53-cleaner.svg?style=flat-square)](https://travis-ci.org/caarlos0/route53-cleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/route53-cleaner?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/route53-cleaner)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/caarlos0)

## Install

### Via homebrew (macOS only):

```console
$ brew install caarlos0/tap/route53-cleaner
```

### Via go get:

```console
$ go get github.com/caarlos0/route53-cleaner
```

### Manually

Download the [latest release](https://github.com/caarlos0/route53-cleaner/releases),
extract it and execute the `route53-cleaner` binary.

## Usage

Just running it will show a list of records you may remove:

```console
$ route53-cleaner
```

## How it works

Route53 Cleaner scans [several resources from your account](/issues/1) and check
your records against those resources addresses, compiling a list of records that 
might not be used.

These records are then printed to the user in an easy-to-pipe format.

Please note that Route53 Cleaner **will never change anything in your account**. You
can check the code or give it read-only keys if you do not trust that affirmation 
(I won't blame you for that).

## Auth

Either by having a `~/.aws/credentials`, `~/.aws/config` or the `AWS_ACCESS_KEY_ID` and 
`AWS_SECRET_ACCESS_KEY` environment variables exported.

More info can be found in the [aws-sdk-go documentation](https://github.com/aws/aws-sdk-go#configuring-credentials).

## Contributing

Please refer to our [contributing guidelines](CONTRIBUTING.md).
