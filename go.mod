module github.com/krakend/krakend-ce/v3

go 1.25.8

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/go-contrib/uuid v1.2.0
	github.com/krakend/bloomfilter/v3 v3.0.0-20260730110053-20a022f5e962
	github.com/krakend/krakend-amqp/v3 v3.0.0-20260731080930-5d5e8c292aeb
	github.com/krakend/krakend-audit/v2 v2.0.0-20260731075904-17ad8e6e0781
	github.com/krakend/krakend-botdetector/v3 v3.0.0-20260730105824-f4b00c30b836
	github.com/krakend/krakend-cel/v3 v3.0.0-20260730105559-6ba9ea59c881
	github.com/krakend/krakend-circuitbreaker/v4 v4.0.0-20260730111621-eec18d30c12d
	github.com/krakend/krakend-cobra/v3 v3.0.0-20260731080528-0ab47e9e709f
	github.com/krakend/krakend-cors/v3 v3.0.0-20260731073552-2729860d05a4
	github.com/krakend/krakend-flexibleconfig/v3 v3.0.0-20260731070931-67f0d356de0c
	github.com/krakend/krakend-gelf/v3 v3.0.0-20260731070446-1470a8946daf
	github.com/krakend/krakend-gologging/v3 v3.0.0-20260730160554-30abc7118c3c
	github.com/krakend/krakend-httpcache/v3 v3.0.0-20260730190833-a6ec1e16f90a
	github.com/krakend/krakend-httpsecure/v3 v3.0.0-20260730185919-03390d6bd5e8
	github.com/krakend/krakend-jose/v3 v3.0.0-20260730185228-6aa13f11f17b
	github.com/krakend/krakend-jsonschema/v3 v3.0.0-20260730183905-5b20c84c62e5
	github.com/krakend/krakend-koanf/v2 v2.0.0-20260818144858-aeb7e1d77767
	github.com/krakend/krakend-lambda/v3 v3.0.0-20260730183553-44bfcf6a660a
	github.com/krakend/krakend-logstash/v3 v3.0.0-20260730182207-3479d8d7f9c5
	github.com/krakend/krakend-lua/v3 v3.0.0-20260730162142-8034a49214ce
	github.com/krakend/krakend-martian/v3 v3.0.0-20260730161248-3c18d144ec94
	github.com/krakend/krakend-metrics/v3 v3.0.0-20260730155513-48daa58645b2
	github.com/krakend/krakend-oauth2-clientcredentials/v3 v3.0.0-20260730152129-ffc709b4c81e
	github.com/krakend/krakend-otel/v2 v2.0.0-20260818145847-dd72b3d29114
	github.com/krakend/krakend-pubsub/v3 v3.0.0-20260730143945-c1a814728f2d
	github.com/krakend/krakend-ratelimit/v4 v4.0.0-20260730143053-a3ce2811a646
	github.com/krakend/krakend-rss/v3 v3.0.0-20260730141845-a9ca05b1a928
	github.com/krakend/krakend-usage/v2 v2.1.0
	github.com/krakend/krakend-xml/v3 v3.0.0-20260730105932-cd80994f0518
	github.com/luraproject/lura/v3 v3.0.0-20260729144624-4b3057d09348
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2
	github.com/spf13/cobra v1.8.1
	golang.org/x/sync v0.22.0
)

require (
	cel.dev/expr v0.25.2 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.22.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.12.0 // indirect
	cloud.google.com/go/kms v1.33.0 // indirect
	cloud.google.com/go/longrunning v1.2.0 // indirect
	cloud.google.com/go/pubsub v1.51.0 // indirect
	cloud.google.com/go/pubsub/v2 v2.6.1 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.2.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.22.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.14.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus v1.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal v1.2.0 // indirect
	github.com/Azure/go-amqp v1.7.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.8.0 // indirect
	github.com/IBM/sarama v1.60.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aws/aws-sdk-go v1.55.8 // indirect
	github.com/aws/aws-sdk-go-v2 v1.43.2 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.33 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.32 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.33 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.50.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.42.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.46.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.2 // indirect
	github.com/aws/smithy-go v1.27.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/gopkg v0.1.4 // indirect
	github.com/bytedance/sonic v1.15.2 // indirect
	github.com/bytedance/sonic/loader v0.5.2 // indirect
	github.com/catalinc/hashcash v0.0.0-20161205220751-e6bc29ff4de9 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/cloudwego/base64x v0.1.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/gin-contrib/sse v1.1.1 // indirect
	github.com/go-jose/go-jose/v3 v3.0.5 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/cel-go v0.29.0 // indirect
	github.com/google/martian v2.1.1-0.20190517191504-25dcb96d9e51+incompatible // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/google/wire v0.7.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.19 // indirect
	github.com/googleapis/gax-go/v2 v2.23.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.2.0 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.7 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-7 // indirect
	github.com/hashicorp/vault/api v1.22.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/parsers/json v1.0.1 // indirect
	github.com/knadh/koanf/parsers/toml v0.1.0 // indirect
	github.com/knadh/koanf/parsers/yaml v1.1.1 // indirect
	github.com/knadh/koanf/providers/env v1.1.0 // indirect
	github.com/knadh/koanf/providers/file v1.2.0 // indirect
	github.com/knadh/koanf/v2 v2.3.6 // indirect
	github.com/krakend/binder v0.0.0-20250826131726-e91a8a754ef8 // indirect
	github.com/krakend/dotnotation v0.2.0 // indirect
	github.com/krakend/flatmap v1.2.0 // indirect
	github.com/krakend/go-auth0/v2 v2.0.4 // indirect
	github.com/krakend/httpcache v1.1.1 // indirect
	github.com/krakend/lru v0.0.0-20250121172718-0e3a6eab620d // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mmcdole/gofeed v1.2.1 // indirect
	github.com/mmcdole/goxpp v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nats.go v1.52.0 // indirect
	github.com/nats-io/nkeys v0.4.16 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.27 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.61.0 // indirect
	github.com/rabbitmq/amqp091-go v1.13.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rs/cors/wrapper/gin v0.0.0-20240830163046-1084d89a1692 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sony/gobreaker/v2 v2.4.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tmthrgd/atomics v0.0.0-20190904060638-dc7a5fcc7e0d // indirect
	github.com/tmthrgd/go-bitset v0.0.0-20190904054048-394d9a556c05 // indirect
	github.com/tmthrgd/go-bitwise v0.0.0-20190904053232-1430ee983fca // indirect
	github.com/tmthrgd/go-byte-test v0.0.0-20190904060354-2794345b9929 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/tmthrgd/go-memset v0.0.0-20190904060434-6fb7a21f88f1 // indirect
	github.com/tmthrgd/go-popcount v0.0.0-20190904054823-afb1ace8b04f // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	github.com/unrolled/secure v1.15.0 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	go.mongodb.org/mongo-driver/v2 v2.8.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.69.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.69.0 // indirect
	go.opentelemetry.io/contrib/propagators/autoprop v0.69.0 // indirect
	go.opentelemetry.io/contrib/propagators/aws v1.44.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.44.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.44.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.44.0 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.66.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	gocloud.dev v0.46.0 // indirect
	gocloud.dev/pubsub/kafkapubsub v0.46.0 // indirect
	gocloud.dev/pubsub/natspubsub v0.46.0 // indirect
	gocloud.dev/pubsub/rabbitpubsub v0.46.0 // indirect
	gocloud.dev/secrets/hashivault v0.45.0 // indirect
	golang.org/x/arch v0.29.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/exp v0.0.0-20260218203240-3dfff04db8fa // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/api v0.291.0 // indirect
	google.golang.org/genproto v0.0.0-20260729162451-8efbd57d26e0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260729162451-8efbd57d26e0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260729162451-8efbd57d26e0 // indirect
	google.golang.org/grpc v1.83.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20191017102106-1550ee647df0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
