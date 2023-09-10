package theoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultAddress = "http://theofficeapi-angelinepinilla.b4a.run"

	defaultRetryWaitMax = 30 * time.Second
)

type Config struct {
	Address string
}

type Client struct {
	baseURL    string
	httpClient *retryablehttp.Client
}

func NewClient(config *Config) (*Client, error) {
	if config.Address == "" {
		config.Address = defaultAddress
	}

	if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address %q: %w", config.Address, err)
	}

	client := retryablehttp.NewClient()
	client.RetryWaitMax = defaultRetryWaitMax
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	return &Client{
		baseURL:    config.Address,
		httpClient: client,
	}, nil
}

type ConnectionsResponse struct {
	Connections []Connection
}

type Connection struct {
	Episode     int    `json:"episode,omitempty"`
	EpisodeName string `json:"episode_name,omitempty"`
	Links       []Link `json:"links,omitempty"`
	Nodes       []Node `json:"nodes,omitempty"`
}

type Link struct {
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
	Value  int    `json:"value,omitempty"`
}

type Node struct {
	ID string `json:"id,omitempty"`
}

type QuotesResponse struct {
	Quotes []Quote
}

type Quote struct {
	Season      int    `json:"season,omitempty"`
	Episode     int    `json:"episode,omitempty"`
	Scene       int    `json:"scene,omitempty"`
	EpisodeName string `json:"episode_name,omitempty"`
	Character   string `json:"character,omitempty"`
	Quote       string `json:"quote,omitempty"`
}

func (c *Client) GetConnections(ctx context.Context, season int) (*ConnectionsResponse, error) {
	path := fmt.Sprintf("/season/%d/format/connections", season)

	resp := &ConnectionsResponse{}
	err := c.do(ctx, "GET", path, nil, &resp.Connections)
	return resp, err
}

func (c *Client) GetQuotes(ctx context.Context, season, episode int) (*QuotesResponse, error) {
	path := fmt.Sprintf("/season/%d", season)
	if episode > 0 {
		path += fmt.Sprintf("/episode/%d", episode)
	} else {
		path += "/format/quotes"
	}
	resp := &QuotesResponse{}
	err := c.do(ctx, "GET", path, nil, &resp.Quotes)
	return resp, err
}

func (c *Client) do(ctx context.Context, method, path string, rq, resp any) error {
	f := func() error {
		logger := hclog.FromContext(ctx).Named("theoffice_client")
		ctx = hclog.WithContext(ctx, logger)
		url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
		var body io.Reader
		if rq != nil {
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(rq); err != nil {
				return fmt.Errorf("encoding request: %w", err)
			}
			body = &buf
		}
		req, err := retryablehttp.NewRequest(method, url, body)
		if err != nil {
			return fmt.Errorf("constructing http request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		logger.Debug("making http request", "method", method, "url", url)
		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		ok := res.StatusCode >= 200 && res.StatusCode < 300
		if !ok {
			resBody, err := io.ReadAll(res.Body)
			if err != nil || string(resBody) == "" {
				return fmt.Errorf("%s %s: bad status (%d)", method, url, res.StatusCode)
			}
			return fmt.Errorf("%s %s: bad status (%d)\n%s", method, url, res.StatusCode, string(resBody))
		}

		return json.NewDecoder(res.Body).Decode(resp)
	}
	return f()
}
