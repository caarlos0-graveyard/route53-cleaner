# Route53 Cleaner

Suggests records that could be deleted from your AWS Route53 hosted zones.

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

## Contributing

Please refer to our [contributing guidelines](CONTRIBUTING.md).
