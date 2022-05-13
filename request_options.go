package http_client

import "net/http"

type CallRequestOption interface {
	before(option *requestOption) error
	after(option *requestOption, response *http.Response)
}

type requestOption struct {
	contentType string
	path        string
	headers     map[string]interface{}
}

func ContentType(contentType string) CallRequestOption {
	return ContentTypeRequestOption{ContentType: contentType}
}

type ContentTypeRequestOption struct {
	ContentType string
}

func (o ContentTypeRequestOption) before(c *requestOption) error {
	c.contentType = o.ContentType
	return nil
}

func (o ContentTypeRequestOption) after(*requestOption, *http.Response) {

}