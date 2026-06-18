package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	ProductionBaseURL = "https://api.mfkessai.co.jp/v2"
	SandboxBaseURL    = "https://sandbox-api.mfkessai.co.jp/v2"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	DryRun     bool
}

func NewClient(sandbox bool) (*Client, error) {
	baseURL := ProductionBaseURL
	if sandbox {
		baseURL = SandboxBaseURL
	}

	apiKey := resolveAPIKey(sandbox)
	if apiKey == "" {
		if sandbox {
			return nil, fmt.Errorf("MFK_SANDBOX_API_KEY or MFK_API_KEY environment variable is required for sandbox")
		}
		return nil, fmt.Errorf("MFK_API_KEY environment variable is required")
	}

	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}, nil
}

func IsSandbox(flagSandbox bool) (bool, error) {
	if flagSandbox {
		return true, nil
	}
	env := os.Getenv("MFK_ENV")
	if env == "" || env == "production" {
		return false, nil
	}
	if env == "sandbox" {
		return true, nil
	}
	return false, fmt.Errorf("invalid MFK_ENV value: %q (must be 'production' or 'sandbox')", env)
}

func resolveAPIKey(sandbox bool) string {
	if sandbox {
		if key := os.Getenv("MFK_SANDBOX_API_KEY"); key != "" {
			return key
		}
	}
	return os.Getenv("MFK_API_KEY")
}

func (c *Client) do(method, path string, query url.Values, body map[string]interface{}) (json.RawMessage, error) {
	u := c.BaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reqBody io.Reader
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	if c.DryRun {
		info := map[string]interface{}{
			"method": method,
			"url":    u,
		}
		if body != nil {
			info["body"] = body
		}
		out, _ := json.MarshalIndent(info, "", "  ")
		fmt.Fprintln(os.Stderr, string(out))
		return nil, nil
	}

	req, err := http.NewRequest(method, u, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("apikey", c.APIKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

func (c *Client) Get(path string, query url.Values) (json.RawMessage, error) {
	return c.do(http.MethodGet, path, query, nil)
}

func (c *Client) Post(path string, body map[string]interface{}) (json.RawMessage, error) {
	return c.do(http.MethodPost, path, nil, body)
}

func (c *Client) Patch(path string, body map[string]interface{}) (json.RawMessage, error) {
	return c.do(http.MethodPatch, path, nil, body)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) GetAllPages(path string, query url.Values) (json.RawMessage, error) {
	var allItems []json.RawMessage

	for {
		raw, err := c.Get(path, query)
		if err != nil {
			return nil, err
		}
		if raw == nil {
			return nil, nil
		}

		var page struct {
			Items      []json.RawMessage `json:"items"`
			Pagination struct {
				HasNext bool   `json:"has_next"`
				End     string `json:"end"`
			} `json:"pagination"`
		}
		if err := json.Unmarshal(raw, &page); err != nil {
			return raw, nil
		}

		allItems = append(allItems, page.Items...)

		if !page.Pagination.HasNext {
			break
		}

		// Advance the cursor using the API-provided `pagination.end` (the last
		// resource ID of this page). Fall back to the last item's `id` if absent.
		cursor := page.Pagination.End
		if cursor == "" && len(page.Items) > 0 {
			var lastItem map[string]interface{}
			if err := json.Unmarshal(page.Items[len(page.Items)-1], &lastItem); err == nil {
				if id, ok := lastItem["id"].(string); ok {
					cursor = id
				}
			}
		}
		if cursor == "" {
			break
		}
		if query == nil {
			query = url.Values{}
		}
		query.Set("after", cursor)
	}

	result := map[string]interface{}{
		"items": allItems,
	}
	return json.Marshal(result)
}

func BuildBody(jsonStr string, overrides map[string]interface{}) (map[string]interface{}, error) {
	body := make(map[string]interface{})

	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &body); err != nil {
			return nil, fmt.Errorf("invalid --json value: %w", err)
		}
	}

	for k, v := range overrides {
		parts := strings.SplitN(k, ".", 2)
		if len(parts) == 2 {
			nested, ok := body[parts[0]].(map[string]interface{})
			if !ok {
				nested = make(map[string]interface{})
			}
			nested[parts[1]] = v
			body[parts[0]] = nested
		} else {
			body[k] = v
		}
	}

	return body, nil
}
