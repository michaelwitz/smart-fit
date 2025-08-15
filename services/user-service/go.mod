module user-service

go 1.24

replace smart-fit/proto => ../../services/db-gateway-service/proto

require (
	github.com/sony/gobreaker v1.0.0
	golang.org/x/crypto v0.41.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)
