package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Get(ctx context.Context, endpoint string, params url.Values) error
	Post(ctx context.Context, endpoint string, params url.Values) error
	Delete(ctx context.Context, endpoint string, params url.Values) error
}

type httpClient struct {
	Host string
	cl   *http.Client
}

func New(host string) HTTPClient {
	cl := http.DefaultClient

	return &httpClient{
		Host: host,
		cl:   cl,
	}
}

func (h *httpClient) Get(ctx context.Context, endpoint string, params url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.buildURL(endpoint, params), nil)
	if err != nil {
		return err
	}
	resp, err := h.cl.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("invalid status code %d, body read error: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("invalid status code %d, body: %s", resp.StatusCode, string(body))
	}

	return err
}

func (h *httpClient) Post(ctx context.Context, endpoint string, params url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.buildURL(endpoint, params), nil)
	if err != nil {
		return err
	}
	resp, err := h.cl.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("invalid status code %d, body read error: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("invalid status code %d, body: %s", resp.StatusCode, string(body))
	}

	return err
}

func (h *httpClient) Delete(ctx context.Context, endpoint string, params url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, h.buildURL(endpoint, params), nil)
	if err != nil {
		return err
	}
	resp, err := h.cl.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("invalid status code %d, body read error: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("invalid status code %d, body: %s", resp.StatusCode, string(body))
	}

	return err
}

func (h *httpClient) buildURL(endpoint string, params url.Values) string {
	return h.Host + "/" + endpoint + "?" + params.Encode()
}
