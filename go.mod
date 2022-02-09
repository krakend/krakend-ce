module github.com/devopsfaith/krakend-ce

go 1.12

require (
	github.com/Azure/go-autorest/autorest v0.11.12 // indirect
	github.com/Unacademy/krakend-error-handler v1.1.1
	github.com/Unacademy/krakend-gin-logger v1.4.2
	github.com/catalinc/hashcash v0.0.0-20161205220751-e6bc29ff4de9 // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/devopsfaith/bloomfilter v1.4.0
	github.com/devopsfaith/flatmap v0.0.0-20200601181759-8521186182fc // indirect
	github.com/devopsfaith/krakend v1.2.0 // indirect
	github.com/devopsfaith/krakend-amqp v1.4.0
	github.com/devopsfaith/krakend-botdetector v1.4.0
	github.com/devopsfaith/krakend-cel v1.4.0
	github.com/devopsfaith/krakend-circuitbreaker v1.4.0
	github.com/devopsfaith/krakend-cobra v1.4.0
	github.com/devopsfaith/krakend-consul v1.4.0
	github.com/devopsfaith/krakend-cors v1.4.0
	github.com/devopsfaith/krakend-etcd v0.0.0-20190425091451-d989a26508d7
	github.com/devopsfaith/krakend-flexibleconfig v1.4.0
	github.com/devopsfaith/krakend-gelf v1.4.0
	github.com/devopsfaith/krakend-gologging v1.4.0
	github.com/devopsfaith/krakend-httpcache v1.4.0
	github.com/devopsfaith/krakend-httpsecure v1.4.0
	github.com/devopsfaith/krakend-influx v1.4.0
	github.com/devopsfaith/krakend-jose v1.4.0
	github.com/devopsfaith/krakend-jsonschema v1.4.0
	github.com/devopsfaith/krakend-lambda v1.4.0
	github.com/devopsfaith/krakend-logstash v1.4.0
	github.com/devopsfaith/krakend-lua v1.4.0
	github.com/devopsfaith/krakend-martian v1.4.0
	github.com/devopsfaith/krakend-metrics v1.4.0
	github.com/devopsfaith/krakend-oauth2-clientcredentials v1.4.0
	github.com/devopsfaith/krakend-opencensus v1.4.1
	github.com/devopsfaith/krakend-pubsub v1.4.0
	github.com/devopsfaith/krakend-ratelimit v1.4.0
	github.com/devopsfaith/krakend-rss v1.4.0
	github.com/devopsfaith/krakend-usage v1.4.0
	github.com/devopsfaith/krakend-viper v1.4.0
	github.com/devopsfaith/krakend-xml v1.4.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-contrib/uuid v1.2.0
	github.com/google/btree v1.0.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/consul v1.6.10 // indirect
	github.com/hashicorp/vault v1.6.0 // indirect
	github.com/influxdata/influxdb v1.7.4 // indirect
	github.com/influxdata/platform v0.0.0-20190117200541-d500d3cf5589 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/kpacha/opencensus-influxdb v0.0.0-20181102202715-663e2683a27c // indirect
	github.com/letgoapp/krakend-consul v0.0.0-20190130102841-7623a4da32a1 // indirect
	github.com/luraproject/lura v1.4.1
	github.com/newrelic/go-agent v3.15.2+incompatible // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/tmthrgd/atomics v0.0.0-20180217065130-6910de195248 // indirect
	github.com/tmthrgd/go-bitwise v0.0.0-20170218093117-01bef038b6bd // indirect
	github.com/tmthrgd/go-byte-test v0.0.0-20170223110042-2eb5216b83f7 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20180828131331-d1fb3dbb16a1 // indirect
	github.com/tmthrgd/go-memset v0.0.0-20180828131805-6f4e59bf1e1d // indirect
	github.com/tmthrgd/go-popcount v0.0.0-20180111143836-3918361d3e97 // indirect
	github.com/unacademy/krakend-auth v1.0.0
	github.com/unacademy/krakend-newrelic v1.0.0-dev.2
	github.com/xeipuuv/gojsonschema v1.2.1-0.20200424115421-065759f9c3d7 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	gocloud.dev v0.21.0 // indirect
	gocloud.dev/pubsub/kafkapubsub v0.21.0 // indirect
	gocloud.dev/pubsub/natspubsub v0.21.0 // indirect
	gocloud.dev/pubsub/rabbitpubsub v0.21.0 // indirect
	gocloud.dev/secrets/hashivault v0.21.0 // indirect
	k8s.io/api v0.20.2 // indirect
)

replace github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 => github.com/m4ns0ur/httpcache v0.0.0-20200426190423-1040e2e8823f

replace github.com/luraproject/lura v1.4.1 => github.com/Unacademy/krakend v1.4.1
