![Krakend logo](https://raw.githubusercontent.com/devopsfaith/krakend.io/master/images/logo.png)

# KrakenD
KrakenD is an extensible, ultra-high performance API Gateway that helps you effortlessly adopt microservices and secure communications. KrakenD is easy to operate and run and scales out without a single point of failure.

**KrakenD Community Edition** (or *KrakenD-CE*) is the open-source distribution of [KrakenD](https://www.krakend.io).

[KrakenD Site](https://www.krakend.io/) | [Documentation](https://www.krakend.io/docs/overview/) | [Blog](https://www.krakend.io/blog/) | [Twitter](https://twitter.com/krakend_io) | [Downloads](https://www.krakend.io/download/)

## Benefits

- **Easy integration** of an ultra-high performance gateway.
- **Effortlessly transition to microservices** and Backend For Frontend implementations.
- **True linear scalability**: Thanks to its **stateless design**, every KrakenD node can operate independently in the cluster without any coordination or centralized persistence.
- **Low operational cost**: +70K reqs/s on a single instance of regular size. Super low memory consumption with high traffic (usually under 50MB w/ +1000 concurrent). Fewer machines. Smaller machines. Lower budget.
- **Platform-agnostic**. Whether you work in a Cloud-native environment (e.g., Kubernetes) or self-hosted on-premises.
- **No vendor lock-in**: Reuse the best existing open-source and proprietary tools rather than having everything in the gateway (telemetry, identity providers, etc.)
- **API Lifecycle**: Using **GitOps** and **declarative configuration**.
- **Decouple clients** from existing services. Create new APIs without changing your existing API contracts.

## Technical features

- **Content aggregation**, composition, and filtering: Create views and mashups of aggregated content from your APIs.
- **Content Manipulation and format transformation**: Change responses, convert transparently from XML to JSON, and vice-versa.
- **Security**: Zero-trust policy, CORS, OAuth, JWT, HSTS, clickjacking protection, HPKP, MIME-Sniffing prevention, XSS protection...
- **Concurrent calls**: Serve content faster than consuming backends directly.
- **SSL** and  **HTTP2** ready
- **Throttling**: Limits of usage in the router and proxy layers
- **Multi-layer rate-limiting** for the end-user and between KrakenD and your services, including bursting, load balancing, and circuit breaker.
- **Telemetry** and dashboards of all sorts: Datadog, Zipkin, Jaeger, Prometheus, Grafana...
- **Extensible** with Go plugins, Lua scripts, Martian, or Google CEL spec.

See the [website](https://www.krakend.io) for more information.

## Download
KrakenD is [packaged and distributed in several formats](https://www.krakend.io/download/). You don't need to clone this repo to use KrakenD unless you want to tweak and build the binary yourself.

## Run
In its simplest form with Docker:

    docker run -it -p "8080:8080" devopsfaith/krakend

Now see [http://localhost:8080/__health](http://localhost:8080/__health). The gateway is listening. Now *CTRL-C* and replace  `/etc/krakend/krakend.json` with your [first configuration](https://designer.krakend.io).

## Build
See the required Go version in the `Makefile`, and then:
```
make build
```

Or, if you don't have or don't want to install `go`, you can build it using the golang docker container:

```
make build_on_docker
```