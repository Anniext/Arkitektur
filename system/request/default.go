package request

import (
	"io"
	"net/http"
)

type HttpRequest struct {
	url    string
	method string
	header map[string]string
	body   io.Reader
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{
		header: make(map[string]string),
	}
}

func (h *HttpRequest) PushHeader(key, value string) *HttpRequest {
	if exit := h.header[key]; exit == "" {
		h.header[key] = value
	}
	return h
}

func (h *HttpRequest) PushMethod(method string) *HttpRequest {
	h.method = method
	return h
}

func (h *HttpRequest) PushUrl(url string) *HttpRequest {
	h.url = url
	return h

}

func (h *HttpRequest) PushBody(body io.Reader) *HttpRequest {
	h.body = body
	return h
}

func (h *HttpRequest) Build() []byte {
	req, _ := http.NewRequest(h.method, h.url, h.body)
	for key, value := range h.header {
		req.Header.Add(key, value)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return body
}
