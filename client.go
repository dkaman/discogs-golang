package discogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dkaman/discogs-golang/internal/options"

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
	Database   *DatabaseService
}

type service struct {
	client *Client
}

func New(opts ...options.Option[Client]) (*Client, error) {
	u, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("default base url failed to parse, this shouldn't happen ever lol. %w", err)
	}

	c := &Client{
		baseURL:     u,
		client:      &http.Client{},
		userAgent:   defaultUserAgent,
		rateLimiter: rate.NewLimiter(rate.Limit(25.0/60.0), 1),
	}

	c.common.client = c
	c.Collection = (*CollectionService)(&c.common)
	c.Identity = (*IdentityService)(&c.common)
	c.Database = (*DatabaseService)(&c.common)

	err = options.Apply(c, opts...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func WithToken(token string) options.Option[Client] {
	return func(c *Client) error {
		// just writing it this way to show 60 request per minute
		c.rateLimiter = rate.NewLimiter(rate.Limit(55.0/60.0), 1)

		auth := fmt.Sprintf("Discogs token=%s", token)
		c.bearerToken = &auth

		_, err := c.Identity.Get()
		if err != nil {
			return fmt.Errorf("error getting identity: %w", err)
		}

		return nil
	}
}

func WithHTTPClient(client *http.Client) options.Option[Client] {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

func (c *Client) NewRequest(method string, urlStr string, body any, opts ...options.Option[http.Request]) (*http.Request, error) {
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

	err = options.Apply(req, opts...)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for rate limiter in http request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	return response, nil
}
