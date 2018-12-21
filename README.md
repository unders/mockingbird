# Mockingbird

Mockingbird is a system testing framework built for the serverless world.

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
GET  http://localhost:8080/ -> redirect-to: http://localhost:8080/v1/dashboard
GET  http://localhost:8080/v1/dashboard

POST http://localhost:8080/v1/tests/
GET  http://localhost:8080/v1/tests/{ID}
GET  http://localhost:8080/v1/tests/

GET  http://localhost:8080/v1/tests/?service=<service>
POST http://localhost:8080/v1/tests/-/services/<service>
```

#e6cece
https://hatchful.shopify.com/editor/customize-logo
https://realfavicongenerator.net/favicon_result?file_id=p1cv6mgrreivv10gm1bu81orr1ln76#.XBwCyhPYoWo

<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="manifest" href="/site.webmanifest.json">
<link rel="mask-icon" href="/safari-pinned-tab.svg" color="#382b2b">
<meta name="apple-mobile-web-app-title" content="Mockingbird">
<meta name="application-name" content="Mockingbird">
<meta name="msapplication-TileColor" content="#da532c">
<meta name="theme-color" content="#e6cece">

<link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png">
<link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png">
<link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png">
<link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png">
<link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png">
<link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png">
<link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png">
<link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png">
<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
<link rel="icon" type="image/png" sizes="192x192"  href="/android-icon-192x192.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="manifest" href="/manifest.json">

## Opencensus
* [guide for client](https://opencensus.io/guides/http/go/net_http/client/)
* [guide for server](https://opencensus.io/guides/http/go/net_http/server/)
* [Go tracing](https://opencensus.io/quickstart/go/tracing/)
* [xray exporter](https://opencensus.io/exporters/supported-exporters/go/xray/)
* [stackdriver exporter](https://opencensus.io/exporters/supported-exporters/go/stackdriver/)

## References
* [rest and long running jobs](https://farazdagi.com/2014/rest-and-long-running-jobs/)
