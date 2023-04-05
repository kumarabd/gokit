package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GET  string = "GET"
	POST        = "POST"
)

type Options struct {
	Type    string
	URL     string
	Headers map[string][]string
	Params  map[string][]string
}

type Handler struct {
	client *http.Client
	req    *http.Request
}

type Response struct {
	Code   int
	Status string
	Data   []byte
}

func New(opts Options) (*Handler, error) {
	//tr := &http.Transport{
	//	MaxIdleConns:       10,
	//	IdleConnTimeout:    30 * time.Second,
	//	DisableCompression: true,
	//}

	client := http.DefaultClient

	url := opts.URL
	if opts.Type == GET {
		// Add Query params
		url = fmt.Sprintf("%s?", opts.URL)
		for key, val := range opts.Params {
			url = fmt.Sprintf("%s&%s=", url, key)
			for _, param := range val {
				url = fmt.Sprintf("%s,", param)
			}
		}
	}

	req, err := http.NewRequest(opts.Type, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = opts.Headers

	return &Handler{
		client: client,
		req:    req,
	}, nil
}

func (h *Handler) Do() (*Response, error) {
	res, err := h.client.Do(h.req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Code:   res.StatusCode,
		Status: res.Status,
		Data:   data,
	}, nil
}
