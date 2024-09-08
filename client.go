package discogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

const (
	version          = "v0.0.1"
	defaultBaseURL   = "https://api.discogs.com"
	defaultUserAgent = "discogs-golang" + "/" + version
)

type Client struct {
	client      *http.Client
	rateLimiter *rate.Limiter
	bearerToken *string
	baseURL     *url.URL
	userAgent   string

	common service

	Collection *CollectionService
	Identity   *IdentityService
}

type service struct {
	client *Client
}

type clientOption func(*Client) error

func New(opts ...clientOption) (*Client, error) {
	u, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("default base url failed to parse, this shouldn't happen ever lol. %w", err)
	}

	c := &Client{
		baseURL:     u,
		client:      &http.Client{},
		userAgent:   defaultUserAgent,
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Minute), 25),
	}

	c.common.client = c
	c.Collection = (*CollectionService)(&c.common)
	c.Identity = (*IdentityService)(&c.common)

	for idx, o := range opts {
		err := o(c)
		if err != nil {
			return nil, fmt.Errorf("error in client option %d: %w", idx, err)
		}
	}

	return c, nil
}

func WithToken(token string) clientOption {
	auth := fmt.Sprintf("Discogs token=%s", token)

	return func(c *Client) error {
		c.bearerToken = &auth

		_, err := c.Identity.Get()
		if err != nil {
			return fmt.Errorf("error getting identity: %w", err)
		}

		c.rateLimiter = rate.NewLimiter(rate.Every(1*time.Minute), 60)

		return nil
	}
}

func WithHTTPClient(client *http.Client) clientOption {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

type requestOption func(req *http.Request) error

func (c *Client) NewRequest(method string, urlStr string, body interface{}, opts ...requestOption) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.bearerToken != nil {
		req.Header.Set("Authorization", *c.bearerToken)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	for _, opt := range opts {
		err = opt(req)
		if err != nil {
			return nil, fmt.Errorf("error with request option: %w", err)
		}
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	response := newResponse(resp)
	return response, nil
}
