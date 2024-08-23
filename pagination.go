package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

var (
	ErrPageDone  = errors.New("no more pages to iterate")
	ErrNilClient = errors.New("provided a nil client to init pager")
)

// stealing from https://vladimir.varank.in/notes/2022/05/a-real-life-use-case-for-generics-in-go-api-for-client-side-pagination/
type Pager[T any] struct {
	pageInfo pageInfo
	client   *Client
}

func (p *Pager[T]) initialize(r *Response) {
	p.pageInfo = r.Paginator
}

func NewPager[T any](r *Response, client *Client) (*T, *Pager[T], error) {
	if client == nil {
		return nil, nil, ErrNilClient
	}

	pager := &Pager[T]{
		client: client,
	}

	pager.initialize(r)

	var respBody T
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return nil, nil, err
	}

	return &respBody, pager, nil
}

func (p *Pager[T]) Next(ctx context.Context) (*T, error) {
	next := p.pageInfo.URLs.Next

	if next == "" {
		return nil, ErrPageDone
	}

	u, err := url.Parse(next)
	if err != nil {
		return nil, err
	}

	req, err := p.client.NewRequest("GET", u.Path+"?"+u.RawQuery, nil)
	if err != nil {
		return nil, err
	}

	var apiResponse *T

	resp, err := p.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return nil, err
	}

	p.pageInfo = resp.Paginator
	return apiResponse, nil
}

func (*Pager[T]) Prev() ([]T, error) {
	return nil, nil
}
