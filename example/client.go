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

	client := NewExampleServiceHTTPClient(cc)

	_, _ = client.GetService(context.Background(), &GetMessageRequest{
		Name: "hello",
	})

	_, _ = client.ListServices(context.Background(), &ListServicesRequest{})
}
