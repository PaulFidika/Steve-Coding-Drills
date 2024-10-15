module steve-fidika/2024-10-12/client

go 1.23

require (
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.34.2
)

require steve-fidika/2024-10-12/server v0.0.0

replace steve-fidika/2024-10-12/server => ../server

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
)
