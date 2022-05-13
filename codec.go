package http_client

import (
	"context"
	"fmt"
	"github.com/dungps/http-client/internal/encoding"
	"io"
	"net/http"
	"strings"
)

func defaultRequestEncoder(ctx context.Context, contentType string, in interface{}) ([]byte, error) {
	body, err := encoding.GetEncoder(contentSubtype(contentType)).Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, err
}

func defaultResponseDecoder(ctx context.Context, res *http.Response, v interface{}) error {
	defer func() {
		_ = res.Body.Close()
	}()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return codecForResponse(res).Unmarshal(data, v)
}

func defaultErrorDecoder(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer func() {
		_ = res.Body.Close()
	}()
	data, err := io.ReadAll(res.Body)
	if err == nil {
		if err = codecForResponse(res).Unmarshal(data, err); err == nil {
			return err
		}
	}

	return fmt.Errorf("response code %d with %v", res.StatusCode, err)
}

func codecForResponse(r *http.Response) encoding.Encoder {
	codec := encoding.GetEncoder(contentSubtype(r.Header.Get("Content-Type")))
	if codec != nil {
		return codec
	}
	return encoding.GetEncoder("json")
}

func contentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}
