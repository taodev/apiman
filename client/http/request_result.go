package http

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/taodev/apiman/test"
)

type RequestResult struct {
	Header map[string]string `yaml:"header,omitempty"`
	Body   string            `yaml:"body,omitempty"`
}

func (r *RequestResult) SetBody(v any) {
	switch v := v.(type) {
	case string:
		r.Body = v
	case []byte:
		buf := make([]byte, hex.EncodedLen(len(v)))
		n := hex.Encode(buf, v)
		r.Body = string(buf[:n])
	default:
		r.Body = fmt.Sprintf("%v", v)
	}
}

func (r *RequestResult) Copy(request http.Request) {
	r.Header = make(map[string]string)
	for k := range request.Header {
		r.Header[k] = request.Header.Get(k)
	}
}

type ResponseResult struct {
	Header map[string]string `yaml:"header,omitempty"`
	Body   string            `yaml:"body,omitempty"`
}

func (r *ResponseResult) SetBody(v any) {
	switch v := v.(type) {
	case string:
		r.Body = v
	case []byte:
		buf := make([]byte, hex.EncodedLen(len(v)))
		n := hex.Encode(buf, v)
		r.Body = string(buf[:n])
	default:
		r.Body = fmt.Sprintf("%v", v)
	}
}

func (r *ResponseResult) Copy(response http.Response) {
	r.Header = make(map[string]string)
	for k := range response.Header {
		r.Header[k] = response.Header.Get(k)
	}
}

type ApiResult struct {
	Name   string `yaml:"name,omitempty"`
	URL    string `yaml:"url"`
	Method string `yaml:"method"`

	ID string `yaml:"id,omitempty"`

	Time time.Duration `yaml:"time"`

	StatusCode int    `yaml:"status_code,omitempty"`
	Status     string `yaml:"status,omitempty"`

	Request  RequestResult  `yaml:"request"`
	Response ResponseResult `yaml:"response"`

	Error string `yaml:"error,omitempty"`

	Tests test.TestManager `yaml:"tests,omitempty"`
}

func (ar *ApiResult) Pass() bool {
	return ar.Tests.Pass()
}
