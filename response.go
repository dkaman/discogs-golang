package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Response struct {
	*http.Response
	Rate      rateLimitInfo
	Paginator pageInfo
}

type responseOption func(*Response) error

type rateLimitInfo struct {
	Limit     int
	Used      int
	Remaining int
}

type pageInfo struct {
	Page    int      `json:"page"`
	Pages   int      `json:"pages"`
	Items   int      `json:"items"`
	PerPage int      `json:"per_page"`
	URLs    pageURLs `json:"urls"`
}

type pageURLs struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

func newResponse(resp *http.Response) *Response {
	response := &Response{Response: resp}
	response.populateRateLimit()
	response.populatePageInfo()
	return response
}

func (r *Response) populateRateLimit() {
	r.Rate.Limit, _ = strconv.Atoi(r.Header.Get("x-discogs-ratelimit"))
	r.Rate.Used, _ = strconv.Atoi(r.Header.Get("x-discogs-ratelimit-used"))
	r.Rate.Remaining, _ = strconv.Atoi(r.Header.Get("x-discogs-ratelimit-remaining"))
}

func (r *Response) populatePageInfo() {
	var p struct {
		Page pageInfo `json:"pagination"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	err = json.Unmarshal(data, &p)
	if err != nil {
		return
	}

	r.Paginator = p.Page
}
