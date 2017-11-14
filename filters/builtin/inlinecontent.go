package builtin

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/zalando/skipper/args"
	"github.com/zalando/skipper/filters"
)

type inlineContent struct {
	text string
	mime string
}

// Creates a filter spec for the inlineContent() filter.
//
// Usage of the filter:
//
//     * -> status(420) -> inlineContent("Enhance Your Calm") -> <shunt>
//
// Or:
//
//     * -> inlineContent("{\"foo\": 42}", "application/json") -> <shunt>
//
// It accepts two arguments: the content and the optional content type.
// When the content type is not set, it tries to detect it using
// http.DetectContentType.
//
// The filter shunts the request with status code 200.
//
func NewInlineContent() filters.Spec {
	return &inlineContent{}
}

func (c *inlineContent) Name() string { return InlineContentName }

func stringArg(a interface{}) (s string, err error) {
	var ok bool
	s, ok = a.(string)
	if !ok {
		err = filters.ErrInvalidFilterParameters
	}

	return
}

func (c *inlineContent) CreateFilter(a []interface{}) (filters.Filter, error) {
	var f inlineContent
	if err := args.Capture(&f.text, args.Optional(&f.mime), a); err != nil {
		return nil, err
	}

	if f.mime == "" {
		f.mime = http.DetectContentType([]byte(f.text))
	}

	return &f, nil
}

func (c *inlineContent) Request(ctx filters.FilterContext) {
	ctx.Serve(&http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{c.mime},
			"Content-Length": []string{strconv.Itoa(len(c.text))},
		},
		Body: ioutil.NopCloser(bytes.NewBufferString(c.text)),
	})
}

func (c *inlineContent) Response(filters.FilterContext) {}