![Krakend logo](https://github.com/devopsfaith/krakend/blob/master/docs/images/krakend.png)

# KrakenD API Gateway
Ultra-High performance API Gateway with middlewares

[![Build Status](https://travis-ci.org/devopsfaith/krakend-ce.svg?branch=master)](https://travis-ci.org/devopsfaith/krakend-ce)

**KrakenD Community Edition** (or *KrakenD-CE*) is the binary distribution of [KrakenD](http://www.krakend.io).

In this repository, you'll find the utils to build the KrakenD binary yourself. If you want to use KrakenD [download the binary](http://www.krakend.io/download/)

[KrakenD Site](http://www.krakend.io/) | [Documentation](http://www.krakend.io/docs/overview/introduction/) | [Blog](http://www.krakend.io/blog) | [Twitter](https://twitter.com/devopsfaith)

## Features
Some of the features you get with KrakenD are:

- **Ultra-High performance** API Gateway
- **Backend for Frontend**
- **Efficient**: Super low memory consumption with high traffic (usually under 50MB w/ +1000 concurrent) and even lower with low traffic (under 5MB)
- **Easy to use**: Declaration of endpoints by just writing the `/url/patterns/and/{variables}`
- **Unlimited** number of backends and endpoints associated with each endpoint. The limit is your kernel.
- **Aggregation** of all the data in the backends for a single endpoint
- **Response composition** and data manipulation (capture, groups, renames...)
- **Response filtering** (whitelist and blacklist)
- **Concurrent** API calls to the backend for the same request
- **Simple configuration**: All application configuration and behavior declared in a `krakend.json`.
- **Friendly**: No development needed to build your gateway. Use the Visual API [Designer](http://www.krakend.io/designer/) (UI to generate the `krakend.json`)
- **SSL**
- **HTTP2** ready
- **Circuit breaker** (fail fast and avoid smashing stressed backends)
- **Bursting** on High-load
- **Logging and statistics** of usage
- **API with statistics**
- **Service Discovery**: DNS SRV, ETCD or add your own
- **Multiple encodings** supported (json, XML, RSS or response as single string)
- **Injections via DSL** in the configuration (Martian)
- **Throttling**: Limits of usage in the router and proxy layers.
- **User quota**: Restrict usage of users by IP or custom headers
- **Basic Firewalling**: Restrict connections by host, drop connections on certain limits
- **Automatic load balancing**
- **HTTP Cache** headers
- **In-memory backend response cache**
- Multiple installation options (bin, docker, rpm, deb, brew)
- **Cloud native**
- **Loved by orchestrators** (Kubernetes, Mesos + Marathon, Nomad, Docker Swarm, and others)
- **Secure:**
    - Support for SSL
    - OAuth client grant
    - JSON Web Tokens (JWT) and JSON Object Signing and Encryption (JOSE)
    - HTTP Strict Transport Security (HSTS)
    - Clickjacking protection
    - HTTP Public Key Pinning (HPKP)
    - MIME-Sniffing prevention
    - Cross-site scripting (XSS) protection
    - Cross-origin resource sharing (CORS)


For a more nice description of the features have a look in the [website](http://www.krakend.io/features/).
## Gateway documentation

Full, comprehensive documentation is viewable on the KrakenD website:

http://www.krakend.io/docs/overview/introduction/

## Build Requirements

- golang 1.11

## Build
```
    make build
```

## Building with docker

If you don't have or don't want to install `go` you can build it using the golang docker container.
```
make docker_build
```

## FPM
You can set up your fpm docker image to run setting `DOCKER_FPM` on the `Makefile`.


## Using the generated packages
The package creates a krakend user to run the service and configures the service to run under systemd.

## Linux Distributions
* just any Linux (using the `tar.gz`)
* ubuntu 14.04, 16.04 (should run un 17.04/10 too)
* Debian 7, 8, 9
* centos 6, 7

```
deb: ubuntu debian
rpm: el6 el7

ubuntu: ubuntu-trusty ubuntu-xenial
debian: debian-wheezy debian-jessie debian-stretch
```
