# Mockingbird

Mockingbird is a system testing framework built for the serverless world.

### Bootstrap

#### Prerequisite
`Go` and `Git` must be install.

#### Install mage
Run command:
```
./tool/install-mage.sh
```
It installs `mage`.

### API

```
GET  http://localhost:8080/ -> redirect-to: http://localhost:8080/v1/dashboard
GET  http://localhost:8080/v1/dashboard

POST http://localhost:8080/v1/tests/
GET  http://localhost:8080/v1/tests/{ID}
GET  http://localhost:8080/v1/tests/

GET  http://localhost:8080/v1/tests/?service=<service>
POST http://localhost:8080/v1/tests/-/services/<service>
```


### Opencensus
* [guide for client](https://opencensus.io/guides/http/go/net_http/client/)
* [guide for server](https://opencensus.io/guides/http/go/net_http/server/)
* [Go tracing](https://opencensus.io/quickstart/go/tracing/)
* [xray exporter](https://opencensus.io/exporters/supported-exporters/go/xray/)
* [stackdriver exporter](https://opencensus.io/exporters/supported-exporters/go/stackdriver/)

### References
* [rest and long running jobs](https://farazdagi.com/2014/rest-and-long-running-jobs/)
