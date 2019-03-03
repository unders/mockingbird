# Mockingbird

Mockingbird is a system testing framework built for the serverless world.

## TODO
* Use S3 Bucket as the database: store result as Json: bucket/{ stats.json | log-uuid.json }
* Add tests for favicons handler
* HTTPS

## Setup

### Prerequisite
 * `Go`
 * `Git`

### Install mage build tool
```
./tool/mage-install.sh
```

## Usage

```
cd tool
mage
```

## API

```
GET  http://localhost:8080/ -> redirect-to: http://localhost:8080/dashboard
GET  http://localhost:8080/dashboard

GET  http://localhost:8080/tests/
GET  http://localhost:8080/tests/{ID}
GET  http://localhost:8080/tests/-/suites/


POST http://localhost:8080/tests/
```


## Opencensus
* [guide for client](https://opencensus.io/guides/http/go/net_http/client/)
* [guide for server](https://opencensus.io/guides/http/go/net_http/server/)
* [Go tracing](https://opencensus.io/quickstart/go/tracing/)
* [xray exporter](https://opencensus.io/exporters/supported-exporters/go/xray/)
* [stackdriver exporter](https://opencensus.io/exporters/supported-exporters/go/stackdriver/)


## References
* [stats](https://github.com/montanaflynn/stats)
* [ulid](https://github.com/ulid/spec)
* [oklog/ulid](https://github.com/oklog/ulid)
* [rest and long running jobs](https://farazdagi.com/2014/rest-and-long-running-jobs/)
* [uikit css](https://getuikit.com/docs/introduction)
* [S3](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html)

