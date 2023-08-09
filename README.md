[![Go](https://github.com/rainu/r-ray/actions/workflows/build.yml/badge.svg)](https://github.com/rainu/r-ray/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/rainu/r-ray/branch/main/graph/badge.svg)](https://codecov.io/gh/rainu/r-ray)
[![Go Report Card](https://goreportcard.com/badge/github.com/rainu/r-ray)](https://goreportcard.com/report/github.com/rainu/r-ray)
[![Go Reference](https://pkg.go.dev/badge/github.com/rainu/r-ray.svg)](https://pkg.go.dev/github.com/rainu/r-ray)

# r-ray

A lightweight proxy application where the **client** can define the request to the target.

# Why?

My main reason is to be able to build powerful PWA. An PWA is an application which is running inside a browser.
Therefore,
this kind of application are bound to the restrictions of the browser. One restriction is the CORS-Policy: The
application
can not make a request to a different origin than the application is hosted. For that issue you can use for example
[CORS Anywhere](https://github.com/Rob--W/cors-anywhere)-Proxy to let set the required CORS-Headers.
This will resolve a lot of issues but not all. If you want to build an application which crawl information from
other origins it can be happened that you have click accept a cookie-banner to get the information. This is not an CORS-
Issue but a Cookie-Issue. And because of browser restrictions an js-application is not able to store/read cookies from
foreign origins!

Don't get me wrong: these are all good mechanism to prevent users from hacking... But in my case this drives me to
implementing
such applications not as PWA but as Desktop-Apps - and it makes me sag :(

# Get the application

You can build it on your own (you will need [golang](https://golang.org/) installed):

```bash
go build -a -installsuffix cgo -o r-ray ./cmd/app
```

Or you can download the release binaries: [here](https://github.com/rainu/r-ray/releases/latest)

Or you can start this application as docker:

```bash
docker pull ghcr.io/rainu/r-ray:main
```

# Usage

Start the application:

```bash
# as binary
CREDENTIALS=user:secret ./r-ray

# as docker container
docker run -p 8080:8080 -e "CREDENTIALS=user:secret" ghcr.io/rainu/r-ray:main
```

After that you can make any request:

```bash
# url -> query encoded url. For example: https://github.com/rainu/r-ray
curl -v -u "user:secret" localhost:8080/?url=https%3A%2F%2Fgithub.com%2Frainu%2Fr-ray
```

# Documentation

## Configuration

| Environment Variable    | Default value             | Description                                                                                                                                         |                                                                                
|-------------------------|---------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|                                                                                
| DEBUG                   | false                     | Debug mode - more logs                                                                                                                              |                                                                                
| BINDING_ADDRESS         | :8080                     | Specifies the TCP address for the server to listen on, in the form "host:port"                                                                      |                                                                                
| CREDENTIALS             |                           | Comma seperated list of credentials in the form "user:password"                                                                                     |
| REQUEST_HEADER_PREFIX   | R-                        | All request header with that prefix will be transfer to target (without prefix). All response Header from target will be prefixed with that prefix. |
| FORWARD_REQUEST_HEADER  | R-Forward-Request-Header  | The header name of the forward request header expressions. See below.                                                                               |
| FORWARD_RESPONSE_HEADER | R-Forward-Response-Header | The header name of the forward response header expressions. See below.                                                                              |
| CORS_ALLOW_ORIGIN       |                           | Comma seperated list of allowed origins. If empty: **every** origin is allowed!                                                                     |
| CORS_ALLOW_METHODS      |                           | Comma seperated list of allowed methods. If empty: **every** method is allowed!                                                                     |
| CORS_ALLOW_HEADERS      |                           | Comma seperated list of allowed headers. If empty: **every** header is allowed!                                                                     |
| CORS_ALLOW_MAX_AGE      |                           | The value of the CORS-Header "Access-Control-Max-Age". If empty: **no** "Access-Control-Max-Age" header will be send.                               |

## Functionality

```
 +--------+          +-------+         +--------+        
 |        |  --1-->  |       |  --2--> |        |
 | client |          | r-ray |         | target |
 |        |  <--4--  |       |  <--3-- |        |
 +--------+          +-------+         +--------+
```

* The target url must be given in the `url`-query parameter. This query parameter should be encoded correctly! 
* The same http-method will be used for target request as the clients request is.
* All header which have the prefix from _REQUEST_HEADER_PREFIX_ (default **R-**) will be transferred to the url's target (without the prefix).
* If the header from _FORWARD_REQUEST_HEADER_ (default **R-Forward-Request-Header**) is sent, all request headers which match the regular expression, will be transferred to the target too
* All headers from the target will be prefixed with _REQUEST_HEADER_PREFIX_ (default **R-**) and sent to the client
* If the header from _FORWARD_RESPONSE_HEADER_ (default **R-Forward-Response-Header**) was sent, all response headers which match the regular expression, will be transferred as is to the client


```
> POST /?url=http://target.com/ HTTP/1.1
> User-Agent: fancy-client
> R-Accept: money
> R-Forward-Request-Header: user-.*
```

Will results in:
```
> POST http://target.com/ HTTP/1.1
> User-Agent: fancy-client
> Accept: money
```
