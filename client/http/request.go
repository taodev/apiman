package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/taodev/apiman/storage"
	"gopkg.in/yaml.v3"
)

type Request struct {
	Param  map[string]string `yaml:"param"`
	Header map[string]string `yaml:"header"`
	Cookie map[string]string `yaml:"cookie"`
	Body   any               `yaml:"body"`

	IgnoreParent bool `yaml:"ignore_parent"`
}

type Response struct {
	Header map[string]string `yaml:"header"`
	Cookie map[string]string `yaml:"cookie"`
	Body   []byte            `yaml:"body"`

	StatusCode int
	Status     string
}

type ApiHttp struct {
	Name     string   `yaml:"name"`
	Method   string   `yaml:"method"`
	URL      string   `yaml:"url"`
	Request  Request  `yaml:"request"`
	Response Response `yaml:"-"`

	BeforeScript       []string `yaml:"before"`
	IgnoreParentBefore bool     `yaml:"ignore_parent_before"`
	AfterScript        []string `yaml:"after"`
	IgnoreParentAfter  bool     `yaml:"ignore_parent_after"`

	// 变量
	Variables map[string]string `yaml:"variables"`
}

func (h *ApiHttp) marshalURL() (fullURL string, err error) {
	rawURL, err := url.Parse(h.URL)
	if err != nil {
		return
	}

	params := url.Values{}
	for k, v := range h.Request.Param {
		params.Set(k, v)
	}

	if len(rawURL.RawQuery) > 0 {
		rawURL.RawQuery += "&" + params.Encode()
	} else {
		rawURL.RawQuery = params.Encode()
	}

	fullURL = rawURL.String()

	return
}

func (h *ApiHttp) processVariables(sessionDB *storage.KeyValueStore) {
	storage.EnvironmentDB.Each(func(k, v string) {
		h.Variables[k] = v
	})

	storage.GlobalDB.Each(func(k, v string) {
		h.Variables[k] = v
	})

	sessionDB.Each(func(k, v string) {
		h.Variables[k] = v
	})
}

func (h *ApiHttp) processBody() (bodyBytes []byte, err error) {
	if h.Request.Body == nil {
		return
	}

	body, ok := h.Request.Body.(string)
	if !ok {
		bodyBytes = h.Request.Body.([]byte)
		return
	}

	t, err := template.New("body").Parse(body)
	if err != nil {
		return
	}

	var writer strings.Builder
	if err = t.Execute(&writer, h.Variables); err != nil {
		return
	}

	out := writer.String()
	h.Request.Body = out
	bodyBytes = []byte(out)

	return
}

func (h *ApiHttp) Do(sessionDB *storage.KeyValueStore) (err error) {
	// 处理环境变量
	h.processVariables(sessionDB)

	var requestBody []byte
	if requestBody, err = h.processBody(); err != nil {
		return
	}

	payload := bytes.NewReader(requestBody)

	rawURL, err := h.marshalURL()
	if err != nil {
		return
	}

	request, err := http.NewRequest(h.Method, rawURL, payload)
	if err != nil {
		return
	}

	url, err := url.Parse(h.URL)
	if err != nil {
		return
	}

	// 设置默认值
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Apiman/1.0.0 (https://apiman.com)")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Host", url.Host)
	request.Header.Add("Connection", "keep-alive")
	request.Header.Set("Accept-Encoding", "gzip")

	for k, v := range h.Request.Header {
		request.Header.Set(k, v)
	}

	for k, v := range h.Request.Cookie {
		request.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var reader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		if reader, err = gzip.NewReader(response.Body); err != nil {
			return
		}
		defer reader.Close()
	} else {
		reader = response.Body
	}

	var body []byte
	if body, err = io.ReadAll(reader); err != nil {
		return
	}

	var resp Response
	resp.Status = response.Status
	resp.StatusCode = response.StatusCode
	resp.Body = body

	resp.Header = make(map[string]string)
	for k := range response.Header {
		resp.Header[k] = response.Header.Get(k)
	}

	resp.Cookie = make(map[string]string)
	cookies := response.Cookies()
	for _, v := range cookies {
		resp.Cookie[v.Name] = v.Value
	}

	h.Response = resp

	return
}

func (h *ApiHttp) Load(workDir string, name string) (err error) {
	// 获取文件绝对路径
	fullpath, err := filepath.Abs(name)
	if err != nil {
		return
	}

	// 获取文件名
	filename := filepath.Base(fullpath)
	fileContent, err := os.ReadFile(fullpath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(fileContent, h)
	if err != nil {
		return
	}

	// 展开body
	if h.Request.Body != nil {
		switch body := h.Request.Body.(type) {
		case string:
		case []byte:
		default:
			var jsonBody []byte
			if jsonBody, err = json.Marshal(body); err != nil {
				return
			}

			h.Request.Body = string(jsonBody)
		}
	}

	// 如果是init.yaml，直接返回
	if filename == "init.yaml" {
		return
	}

	// 读取上级目录配置
	nodes := []*ApiHttp{h}
	dir := filepath.Dir(fullpath)
	for strings.Contains(dir, workDir) {
		initConfig := filepath.Join(dir, "init.yaml")
		dir = filepath.Dir(dir)

		// 配置文件是否存在
		if _, err = os.Stat(initConfig); err != nil {
			err = nil
			continue
		}

		init := new(ApiHttp)
		if err = init.Load(workDir, initConfig); err != nil {
			return
		}

		nodes = append(nodes, init)
	}

	if err = h.expand(nodes); err != nil {
		return
	}

	return
}

func (h *ApiHttp) expand(nodes []*ApiHttp) (err error) {
	// 展开父类设置
	var n1 ApiHttp
	n1.Request.Param = make(map[string]string)
	n1.Request.Header = make(map[string]string)
	n1.Request.Cookie = make(map[string]string)
	n1.Variables = make(map[string]string)

	for i := len(nodes) - 1; i >= 0; i-- {
		if err = h.comboine(&n1, *nodes[i]); err != nil {
			return
		}
	}

	h.Method = n1.Method
	h.URL = n1.URL

	h.Request.Param = n1.Request.Param
	h.Request.Header = n1.Request.Header
	h.Request.Cookie = n1.Request.Cookie

	h.BeforeScript = n1.BeforeScript
	h.AfterScript = n1.AfterScript

	h.Variables = n1.Variables

	return
}

func (h *ApiHttp) comboine(n1 *ApiHttp, n2 ApiHttp) (err error) {
	if len(n2.Method) > 0 {
		n1.Method = n2.Method
	}

	if len(n2.URL) > 0 {
		if strings.Contains(n2.URL, "http://") || strings.Contains(n2.URL, "https://") {
			n1.URL = n2.URL
		} else {
			if n1.URL, err = url.JoinPath(n1.URL, n2.URL); err != nil {
				return
			}
		}
	}

	if n2.Request.IgnoreParent {
		n1.Request = n2.Request
		if n1.Request.Param == nil {
			n1.Request.Param = make(map[string]string)
		}

		if n1.Request.Header == nil {
			n1.Request.Header = make(map[string]string)
		}

		if n1.Request.Cookie == nil {
			n1.Request.Cookie = make(map[string]string)
		}
	} else {
		for k, v := range n2.Request.Param {
			n1.Request.Param[k] = v
		}

		for k, v := range n2.Request.Header {
			n1.Request.Header[k] = v
		}

		for k, v := range n2.Request.Cookie {
			n1.Request.Cookie[k] = v
		}
	}

	if n2.IgnoreParentBefore {
		n1.BeforeScript = n2.BeforeScript
	} else {
		if len(n2.BeforeScript) > 0 {
			n1.BeforeScript = append(n1.BeforeScript, n2.BeforeScript...)
		}
	}

	if n2.IgnoreParentAfter {
		n1.AfterScript = n2.AfterScript
	} else {
		if len(n2.AfterScript) > 0 {
			n1.AfterScript = append(n1.AfterScript, n2.AfterScript...)
		}
	}

	for k, v := range n2.Variables {
		n1.Variables[k] = v
	}

	return
}
