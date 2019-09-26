module github.com/openrm/krakend-ce

go 1.12

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20190108154635-47c0da630f72

require (
	contrib.go.opencensus.io/exporter/jaeger v0.0.0-20190424224017-5b8293c22f36 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.0.0-20190424224027-f02a6e68f94d // indirect
	contrib.go.opencensus.io/exporter/zipkin v0.0.0-20190424224031-c96617f51dc6 // indirect
	github.com/Azure/azure-service-bus-go v0.4.1 // indirect
	github.com/Azure/go-autorest v11.6.0+incompatible // indirect
	github.com/PuerkitoBio/goquery v1.4.0 // indirect
	github.com/andybalholm/cascadia v1.0.0 // indirect
	github.com/auth0-community/go-auth0 v1.0.0
	github.com/catalinc/hashcash v0.0.0-20161205220751-e6bc29ff4de9 // indirect
	github.com/clbanning/mxj v0.0.0-20180418195244-1f00e0bf9bac // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	github.com/coreos/etcd v3.3.15+incompatible // indirect
	github.com/devopsfaith/krakend v0.0.0-20190921111907-6ff3a6860ce0
	github.com/devopsfaith/krakend-ce v0.0.0-20190917135805-07029e9a7b49
	github.com/devopsfaith/krakend-cel v0.0.0-20190502130550-d6872fd4f97e
	github.com/devopsfaith/krakend-circuitbreaker v0.0.0-20190206135831-673caf8e355a
	github.com/devopsfaith/krakend-cobra v0.0.0-20190403130617-3f056666a23e
	github.com/devopsfaith/krakend-consul v0.0.0-20190130102841-7623a4da32a1
	github.com/devopsfaith/krakend-cors v0.0.0-20180906120953-c9b7eb26914f
	github.com/devopsfaith/krakend-etcd v0.0.0-20180215165313-bd645943ff8c
	github.com/devopsfaith/krakend-flexibleconfig v0.0.0-20190408143848-fc4ef2b4d5cf
	github.com/devopsfaith/krakend-gelf v0.0.0-20181019222239-59c0250b1c60
	github.com/devopsfaith/krakend-gologging v0.0.0-20190131142345-f3f256584ecc
	github.com/devopsfaith/krakend-httpcache v0.0.0-20181030153148-8474476ff874
	github.com/devopsfaith/krakend-httpsecure v0.0.0-20180922151646-cce73b27c717
	github.com/devopsfaith/krakend-jsonschema v0.0.0-20190124184701-5705a5015d7a
	github.com/devopsfaith/krakend-logstash v0.0.0-20190131142205-17f4745d3502
	github.com/devopsfaith/krakend-martian v0.0.0-20190424133031-29314a524a91
	github.com/devopsfaith/krakend-metrics v0.0.0-20190114200758-1e2c2a1f6a62
	github.com/devopsfaith/krakend-oauth2-clientcredentials v0.0.0-20190206125733-11a9f7170c44
	github.com/devopsfaith/krakend-opencensus v0.0.0-20190425142549-a584d6fd2cc1
	github.com/devopsfaith/krakend-pubsub v0.0.0-20190424155946-2884ffb54959
	github.com/devopsfaith/krakend-ratelimit v0.0.0-20190404110207-d63774e96e82
	github.com/devopsfaith/krakend-rss v0.0.0-20180408220939-4c18c62a99ee
	github.com/devopsfaith/krakend-usage v0.0.0-20181025134340-476779c0a36c
	github.com/devopsfaith/krakend-viper v0.0.0-20190407170411-1cbb76813774
	github.com/devopsfaith/krakend-xml v0.0.0-20180408220837-5ce94062a4cc
	github.com/geetarista/go-bloomd v0.0.0-20140722181834-7f8e8a358bec
	github.com/getsentry/sentry-go v0.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-contrib/uuid v1.2.0
	github.com/google/cel-go v0.2.0 // indirect
	github.com/influxdata/influxdb v1.7.4 // indirect
	github.com/juju/ratelimit v0.0.0-20171026090426-59fac5042749 // indirect
	github.com/kpacha/opencensus-influxdb v0.0.0-20181102202715-663e2683a27c // indirect
	github.com/letgoapp/krakend-influx v0.0.0-20190214142340-d2fc9466bb3a
	github.com/mmcdole/gofeed v1.0.0-beta2 // indirect
	github.com/mmcdole/goxpp v0.0.0-20170720115402-77e4a51a73ed // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/openrm/krakend-jose v0.0.0-20190925083548-04bd6fcc5643
	github.com/openrm/module-tracing-golang v1.0.9
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/gobreaker v0.0.0-20170530031423-e9556a45379e // indirect
	github.com/spf13/viper v1.3.2 // indirect
	github.com/streadway/amqp v0.0.0-20190402114354-16ed540749f6 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	github.com/unrolled/secure v0.0.0-20171102162350-0f73fc7feba6 // indirect
	go.opencensus.io v0.22.1
	golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392 // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7 // indirect
	gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20180326133423-4dbb9d721348 // indirect
	gopkg.in/gin-contrib/cors.v1 v1.0.0-20170318125340-cf4846e6a636 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
)

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190909230951-414d861bb4ac

replace github.com/hashicorp/vault-plugin-auth-pcf => github.com/hashicorp/vault-plugin-auth-cf v0.0.0-20190821162840-1c2205826fee

replace gopkg.in/gin-contrib/cors.v1 => github.com/gin-contrib/cors v1.3.0

replace gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.22.0

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1

replace github.com/go-xorm/core => xorm.io/core v0.7.0

replace github.com/Unknwon/com => github.com/unknwon/com v1.0.1

replace github.com/Unknwon/i18n => github.com/unknwon/i18n v0.0.0-20190805065654-5c6446a380b6

replace gopkg.in/stretchr/testify.v1 => github.com/stretchr/testify v1.4.0

replace git.apache.org/thrift.git@v0.12.0 => github.com/apache/thrift v0.12.0

replace gopkg.in/src-d/go-git-fixtures.v3 => github.com/src-d/go-git-fixtures v3.5.0+incompatible
