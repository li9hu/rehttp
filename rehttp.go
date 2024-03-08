package rehttp

import (
	"bytes"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/url"
	"time"
)

type HTTP struct {
	Client        *retryablehttp.Client
	ResponseLimit int64
	Timeout       int64
	RetryMax      int
}

var DefaultClient = httpDefault()

// NewDefault 默认初始化的HTTP
func httpDefault() *HTTP {
	var defaultHTTP = &HTTP{
		Client:        initHTTP(10, 3),
		ResponseLimit: 0,
	}
	return defaultHTTP
}

func GlobalHttp(h *HTTP) {
	DefaultClient.Timeout = h.Timeout
	DefaultClient.RetryMax = h.RetryMax
	DefaultClient.ResponseLimit = h.ResponseLimit
	DefaultClient.Client = initHTTP(h.Timeout, h.RetryMax)
}

// New 创建HTTPClient
func New(http *HTTP) *HTTP {
	h := &HTTP{
		Client:        initHTTP(http.Timeout, http.RetryMax),
		ResponseLimit: http.ResponseLimit,
	}
	return h
}

// Post 发送POST请求
func (h *HTTP) Post(urlStr string, headers map[string]string, requestBody []byte) *Result {
	u, err := url.Parse(urlStr)
	if err != nil {
		return &Result{URL: urlStr, Err: err}
	}

	startTime := time.Now()
	req, err := retryablehttp.NewRequest("POST", u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return &Result{URL: u.String(), Err: err}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return &Result{URL: u.String(), Err: err}
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	// 读取响应结果，有长度限制则截取部分内容
	var bodyReader io.Reader = resp.Body
	if h.ResponseLimit != 0 {
		bodyReader = io.LimitReader(resp.Body, h.ResponseLimit)
	}
	b, err := io.ReadAll(bodyReader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return &Result{URL: u.String(), Err: err}
	}

	// 执行花费时间
	durationSeconds := time.Since(startTime).Seconds()

	return &Result{
		URL:          u.String(),
		ResponseBody: string(b),
		StatusCode:   resp.StatusCode,
		Duration:     durationSeconds,
	}
}

func Post(urlStr string, headers map[string]string, requestBody []byte) *Result {
	result := DefaultClient.Post(urlStr, headers, requestBody)
	return result
}

// Get 发送GET请求
func (h *HTTP) Get(t string, headers map[string]string) *Result {
	u, err := url.Parse(t)
	if err != nil {
		return &Result{URL: t, Err: err}
	}

	startTime := time.Now()
	req, err := retryablehttp.NewRequest("GET", u.String(), nil)
	if err != nil {
		return &Result{URL: u.String(), Err: err}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return &Result{URL: u.String(), Err: err}
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	// 读取响应结果，有长度限制则截取部分内容
	var bodyReader io.Reader = resp.Body
	if h.ResponseLimit != 0 {
		bodyReader = io.LimitReader(resp.Body, h.ResponseLimit)
	}
	b, err := io.ReadAll(bodyReader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return &Result{URL: u.String(), Err: err}
	}

	// 执行花费时间
	duration := time.Since(startTime).Seconds()

	return &Result{
		URL:          u.String(),
		ResponseBody: string(b),
		StatusCode:   resp.StatusCode,
		Duration:     duration,
	}
}

func Get(t string, headers map[string]string) *Result {
	result := DefaultClient.Get(t, headers)
	return result
}

// 自带请求重试的HTTPClient
func initHTTP(timeout int64, retryMax int) *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryMax = retryMax
	retryClient.HTTPClient.Timeout = time.Duration(timeout) * time.Second
	return retryClient
}
