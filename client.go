package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
	common      service
	baseURL     *url.URL
	userAgent   string

	Collection *CollectionService
}

type service struct {
	client *Client
}

func NewClient(authToken string, httpClient *http.Client) (*Client, error) {
	var b *string
	if authToken != "" {
		t := fmt.Sprintf("Discogs token=%s", authToken)
		b = &t
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpClient2 := *httpClient

	c := &Client{
		client:      &httpClient2,
		bearerToken: b,
	}

	err := c.initialize()

	return c, err
}

func (c *Client) initialize() error {
	if c.baseURL == nil {
		u, err := url.Parse(defaultBaseURL)
		if err != nil {
			return err
		}
		c.baseURL = u
	}

	if c.userAgent == "" {
		c.userAgent = defaultUserAgent
	}

	if c.bearerToken == nil {
		c.rateLimiter = rate.NewLimiter(rate.Every(1*time.Minute), 25)
	} else {
		c.rateLimiter = rate.NewLimiter(rate.Every(1*time.Minute), 60)
	}

	c.common.client = c
	c.Collection = (*CollectionService)(&c.common)

	return nil
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

func main() {
	ctx := context.Background()
	authToken := os.Getenv("pat")
	c, err := NewClient(authToken, nil)
	if err != nil {
		fmt.Printf("error creating client: %s", err)
		return
	}

	f, err := c.Collection.GetReleasesByFolder(ctx, "dallaskaman", 0)
	if err != nil {
		fmt.Printf("error getting folders: %s", err)
		return
	}

	for _, l := range f {
		fmt.Printf("release %s\n", l.BasicInfo.Title)
	}

	g, err := c.Collection.GetFolderByRelease(ctx, "dallaskaman", 12245977)

	for _, gg := range g {
		fmt.Printf("release - by folder: %s\n", gg.BasicInfo.Title)
	}

	// req, err := c.NewRequest("GET", "artists/1/releases?page=1&per_page=1", nil)
	// if err != nil {
	// 	fmt.Printf("error creating request: %s", err)
	// 	return
	// }
	// reqDump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	fmt.Printf("error dumping request: %s\n", err)
	// 	return
	// }

	// fmt.Printf("REQUEST:\n%s", string(reqDump))

	// resp, err := c.Do(ctx, req)
	// if err != nil {
	// 	fmt.Printf("error sending request: %s", err)
	// 	return
	// }
	// respDump, err := httputil.DumpResponse(resp.Response, true)
	// if err != nil {
	// 	fmt.Printf("error dumping response: %s\n", err)
	// }
	// fmt.Printf("RESPONSE:\n%s", string(respDump))
}
