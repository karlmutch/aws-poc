# Amazon Web Services Proof of Concept and Prototypes

This project has been motivated by a desire to test or examine the character of the AWS golang API, and in addition some of the features of AWS.

When employing tools within this project you will be services that are fee based services offered by Amazon.
Before running any tools within this package be sure you understand the financial implications of doing so.

This project is a personal discovery endeavor and should not be used to in way to imply the software architecture I am using personally, or in my professional capacity.  This project is simply a way for me to experiment with some of AWS APIs to validate and investigate their ease of use and the utility of various features.

## Configuration


The tools within this project have an expectation that the AWS credentials will be made available using the prototypical AWS environment variables AWS_ACCESS_KEY_ID, AWS_DEFAULT_REGION, and AWS_SECRET_ACCESS_KEY.  Beware of the security implications of passing these environment variables around and never store them on github or in any other accessible location.

## Build

The code in the repository can be compiled, and run using the included Makefile.  The Makefile is provided in order that we can ignore `GOPATH` concerns when using this code.

The following make targets are available:

 - `make help` to get help
 - `make` to build the binary (in `bin/`)
 - `make test` to run tests
 - `make test-verbose` to run tests in verbose mode
 - `make test-race` for race tests
 - `make test-xml` for tests with xUnit-compatible output
 - `make test-coverage` for test coverage (will output `index.html`, `coverage.xml` and `profile.out` in `test/coverage.*/`.
 - `make test PKG=helloworld/hello` to restrict test to a package
 - `make clean`
 - `make vendor` to retrieve dependencies
 - `make lint` to run golint
 - `make fmt` to run gofmt

## Usage

When running the make logging with the code base can be controlled using the logxi environment variables described at https://github.com/mgutz/logxi.  For example running `LOGXI=*=INF make test-verbose` will show some limited summary output from the test run.


## License

This example code is distributed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
see LICENSE.txt and NOTICE.txt for more information.


Have fun,
<p>
karlmutch@gmail.com

