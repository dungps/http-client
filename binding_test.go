package http_client

import (
	"github.com/dungps/http-client/internal/testdata"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	url := EncodeURL(
		"http://example.com/helloworld/{name}",
		&testdata.HelloWorld{Name: "hello"},
		false,
	)
	t.Log(url)
	if url != "http://example.com/helloworld/hello" {
		t.Fatalf("proto path not expected!actual: %s", url)
	}

	url = EncodeURL(
		"http://example.com/helloworld/{name}/sub/{sub.name}",
		&testdata.HelloWorld{
			Name: "hello",
			Sub: &testdata.HelloWorldSub{
				Name: "client",
			},
		},
		false,
	)
	t.Log(url)
	if url != "http://example.com/helloworld/hello/sub/client" {
		t.Fatalf("proto path not expected!actual: %s", url)
	}

	url = EncodeURL(
		"http://example.com/helloworld/{name}",
		&testdata.HelloWorld{
			Name: "hello",
			Sub: &testdata.HelloWorldSub{
				Name: "client",
			},
		},
		true,
	)
	t.Log(url)
	if url != "http://example.com/helloworld/hello?sub.name=client" {
		t.Fatalf("proto path not expected!actual: %s", url)
	}
}
