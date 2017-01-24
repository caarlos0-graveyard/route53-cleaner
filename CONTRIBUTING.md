# Contributing

By participating to this project, you agree to abide our [code of
conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`route53-cleaner` is written in [Go](https://golang.org/).

Prerequisites are:

* Build:
  * `make`
  * [Go 1.7+](http://golang.org/doc/install)

Clone `route53-cleaner` from source into `$GOPATH`:

```sh
$ go get github.com/caarlos0/route53-cleaner
$ cd $GOPATH/src/github.com/caarlos0/route53-cleaner
```

Install the build and lint dependencies:

``` sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

``` sh
$ make test
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

``` sh
$ go build
```

When you are satisfied with the changes, we suggest you run:

``` sh
$ make ci
```

Which runs all the linters and tests.

## Submit a pull request

Push your branch to your `route53-cleaner` fork and open a pull request against the
master branch.
