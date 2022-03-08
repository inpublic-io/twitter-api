module github.com/inpublic-io/twitter-api

go 1.17

require (
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/influxdata/influxdb-client-go/v2 v2.8.0
	github.com/inpublic-io/inpublicapis/twitter v0.0.4
	github.com/joho/godotenv v1.4.0
	github.com/vniche/twitter-go v0.0.6
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.8.2 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/cors v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	nhooyr.io/websocket v1.8.6 // indirect
)

replace github.com/nats-io/nats-server/v2 => github.com/nats-io/nats-server/v2 v2.7.2

replace github.com/apache/thrift => github.com/apache/thrift v0.16.0

replace github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1

replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2

replace github.com/miekg/dns => github.com/miekg/dns v1.1.46
