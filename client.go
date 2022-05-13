package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/dungps/http-client/internal/encoding/form"
	_ "github.com/dungps/http-client/internal/encoding/json"
	_ "github.com/dungps/http-client/internal/encoding/xml"
)

type DecodeErrorFunc func(ctx context.Context, res *http.Response) error

type EncodeRequestFunc func(ctx context.Context, contentType string, in interface{}) (body []byte, err error)

type DecodeResponseFunc func(ctx context.Context, res *http.Response, out interface{}) error

type ClientOption func(*clientOptions)

type clientOptions struct {
	tlsConfig    *tls.Config
	timeout      time.Duration
	baseURL      string
	transport    http.RoundTripper
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
}

func WithTLSConfig(tlsConf *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConfig = tlsConf
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

func WithBaseURL(url string) ClientOption {
	return func(o *clientOptions) {
		o.baseURL = url
	}
}

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = transport
	}
}

func WithResponseDecoder(decoder DecodeResponseFunc) ClientOption {
	return func(o *clientOptions) {
		o.decoder = decoder
	}
}

func WithRequestEncoder(encoder EncodeRequestFunc) ClientOption {
	return func(o *clientOptions) {
		o.encoder = encoder
	}
}

func WithErrorDecoder(errDecoder DecodeErrorFunc) ClientOption {
	return func(o *clientOptions) {
		o.errorDecoder = errDecoder
	}
}

func NewClient(opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		timeout:      60 * time.Second,
		transport:    http.DefaultTransport,
		encoder:      defaultRequestEncoder,
		decoder:      defaultResponseDecoder,
		errorDecoder: defaultErrorDecoder,
	}

	for _, o := range opts {
		o(&options)
	}

	insecure := options.tlsConfig == nil
	target, err := parseURL(options.baseURL, insecure)
	if err != nil {
		return nil, err
	}

	return &Client{
		opts:     options,
		insecure: insecure,
		target:   target,
		c: &http.Client{
			Timeout:   options.timeout,
			Transport: options.transport,
		},
	}, nil
}

type Client struct {
	opts     clientOptions
	c        *http.Client
	target   *url.URL
	insecure bool
}

func (c *Client) Request(ctx context.Context, method, path string, in interface{}, out interface{}, opts ...CallRequestOption) error {
	var body io.Reader
	opt := requestOption{
		contentType: "application/json",
		path:        path,
	}
	for _, o := range opts {
		if err := o.before(&opt); err != nil {
			return err
		}
	}

	if in != nil {
		data, err := c.opts.encoder(ctx, opt.contentType, in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	u := fmt.Sprintf("%s://%s%s", c.target.Scheme, c.target.Host, path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return err
	}
	if opt.contentType != "" {
		req.Header.Set("Content-Type", opt.contentType)
	}

	handler := func(ctx context.Context, in interface{}) error {
		res, err := c.do(req.WithContext(ctx))
		if res != nil {
			for _, o := range opts {
				o.after(&opt, res)
			}
		}
		if err != nil {
			return err
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if err := c.opts.decoder(ctx, res, out); err != nil {
			return err
		}

		return nil
	}

	return handler(ctx, req)
}

func (c *Client) Do(req *http.Request, opts ...CallRequestOption) (*http.Response, error) {
	opt := requestOption{
		contentType: "application/json",
		path:        req.URL.Path,
	}
	for _, o := range opts {
		if err := o.before(&opt); err != nil {
			return nil, err
		}
	}

	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.c.Do(req)
	if err == nil {
		err = c.opts.errorDecoder(req.Context(), resp)
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func parseURL(baseUrl string, insecure bool) (*url.URL, error) {
	if !strings.Contains(baseUrl, "://") {
		if insecure {
			baseUrl = "http://" + baseUrl
		} else {
			baseUrl = "https://" + baseUrl
		}
	}
	return url.Parse(baseUrl)
}
