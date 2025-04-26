module github.com/snapser/simplegs

go 1.20

require (
	agones.dev/agones v1.33.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.3.0
	github.com/stretchr/testify v1.8.4
	google.golang.org/genproto/googleapis/api v0.0.0-20230711160842-782d3b101e98
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace agones.dev/agones => ./
