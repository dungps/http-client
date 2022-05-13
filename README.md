# HTTP Client

> This module inspired by Kratos framework

HTTP Client is a framework which easy to create HTTP client with protobuf

## Usage

1. Create proto

```protobuf
syntax = "proto3";

package example;

import "google/api/annotations.proto";

service HelloWorld {
  rpc GetMessage(Request) returns (Response) {
    option(google.api.http) = {
      get: ""
    };
  }
}

message Request {
  string name = 1;
}

message Response {
  string message = 1;
}
```

2. Generate go with protobuf

This is require some cli
- [protoc](https://github.com/protocolbuffers/protobuf-go)
- [protoc-gen-go-client](https://github.com/dungps/http-client/tree/master/cmd/protoc-gen-go-client)

```shell
protoc --go_out=paths=source_relative:. --go_client=paths=source_relative:. example.proto
```

3. Using

```go
package main

import (
	"context"
	http_client "github.com/dungps/http-client"
)

func main() {
	cc, err := http_client.NewClient(
		http_client.WithBaseURL("https://example.com"),
	)
	if err != nil {
		panic(err)
	}

	client := NewHelloWorldHTTPClient(cc)

	_, _ = client.GetMessage(context.Background(), &Request{
		Name: "hello",
	})
}
```