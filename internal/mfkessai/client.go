package mfkessai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.mfkessai.co.jp"
)

// Client is a client for MoneyForward Kessai API
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new MF Kessai API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Billing represents a billing record
type Billing struct {
	Object     string   `json:"object"`
	ID         string   `json:"id"`
	CustomerID string   `json:"customer_id"`
	Amount     int      `json:"amount"`
	IssueDate  string   `json:"issue_date"`
	DueDate    string   `json:"due_date"`
	Status     string   `json:"status"`
	InvoiceIDs []string `json:"invoice_ids"`
}

// BillingsResponse represents the response from /v2/billings
type BillingsResponse struct {
	Object     string `json:"object"`
	Pagination struct {
		Total       int    `json:"total"`
		Limit       int    `json:"limit"`
		HasNext     bool   `json:"has_next"`
		HasPrevious bool   `json:"has_previous"`
		Start       string `json:"start"`
		End         string `json:"end"`
	} `json:"pagination"`
	Items []Billing `json:"items"`
}

// BillingSearchParams contains parameters for searching billings
type BillingSearchParams struct {
	IssueDateFrom string
	IssueDateTo   string
	DueDateFrom   string
	DueDateTo     string
	Status        []string
}

// DownloadSignedURLRequest represents the request to get download signed URL
type DownloadSignedURLRequest struct {
	Type string `json:"type"`
}

// DownloadSignedURLResponse represents the response from download_signed_url endpoint
type DownloadSignedURLResponse struct {
	Object string `json:"object"`
	Items  []struct {
		SignedURL string    `json:"signed_url"`
		ExpiredAt time.Time `json:"expired_at"`
		Type      string    `json:"type"`
	} `json:"items"`
}

// GetBillings retrieves billings within the specified date range
func (c *Client) GetBillings(issueDateFrom, issueDateTo string) ([]Billing, error) {
	return c.GetBillingsWithParams(BillingSearchParams{
		IssueDateFrom: issueDateFrom,
		IssueDateTo:   issueDateTo,
	})
}

// GetBillingsWithParams retrieves billings with custom search parameters
// Uses the /v2/billings/qualified endpoint for invoice system compliance
func (c *Client) GetBillingsWithParams(params BillingSearchParams) ([]Billing, error) {
	var allBillings []Billing
	var startingAfter string

	for {
		// Use /v2/billings/qualified endpoint for invoice system (インボイス制度) compliance
		url := fmt.Sprintf("%s/v2/billings/qualified?limit=100", baseURL)

		// Add date filters
		if params.IssueDateFrom != "" {
			url += fmt.Sprintf("&issue_date_from=%s", params.IssueDateFrom)
		}
		if params.IssueDateTo != "" {
			url += fmt.Sprintf("&issue_date_to=%s", params.IssueDateTo)
		}
		if params.DueDateFrom != "" {
			url += fmt.Sprintf("&due_date_from=%s", params.DueDateFrom)
		}
		if params.DueDateTo != "" {
			url += fmt.Sprintf("&due_date_to=%s", params.DueDateTo)
		}

		// Add status filters
		for _, status := range params.Status {
			url += fmt.Sprintf("&status=%s", status)
		}

		if startingAfter != "" {
			url += fmt.Sprintf("&starting_after=%s", startingAfter)
		}

		// Debug: Log the request URL
		fmt.Printf("DEBUG: Fetching URL: %s\n", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("apikey", c.apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}

		var billingsResp BillingsResponse
		if err := json.NewDecoder(resp.Body).Decode(&billingsResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		// Debug: Log pagination info
		fmt.Printf("DEBUG: Page total=%d, has_next=%v, items_count=%d, end=%s\n",
			billingsResp.Pagination.Total, billingsResp.Pagination.HasNext, len(billingsResp.Items), billingsResp.Pagination.End)

		allBillings = append(allBillings, billingsResp.Items...)

		// Check if there are more pages
		if !billingsResp.Pagination.HasNext {
			break
		}

		// Use the End cursor for pagination
		newStartingAfter := billingsResp.Pagination.End
		if newStartingAfter == startingAfter || newStartingAfter == "" {
			// Prevent infinite loop if API returns same cursor
			fmt.Printf("DEBUG: WARNING - Pagination cursor not advancing (current=%s, new=%s). Breaking loop.\n", startingAfter, newStartingAfter)
			break
		}
		startingAfter = newStartingAfter
	}

	return allBillings, nil
}

// GetDownloadSignedURL gets a signed URL for downloading the invoice PDF
// The endpoint format is: /v2/billings/{billing_id}/issues/i/{invoice_id}/download_signed_url
// where 'i' is a literal path segment indicating invoice type
func (c *Client) GetDownloadSignedURL(billingID, invoiceID string) (string, error) {
	url := fmt.Sprintf("%s/v2/billings/%s/issues/i/%s/download_signed_url", baseURL, billingID, invoiceID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var signedURLResp DownloadSignedURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&signedURLResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Return the first PDF URL from items
	if len(signedURLResp.Items) > 0 {
		for _, item := range signedURLResp.Items {
			if item.Type == "pdf" {
				return item.SignedURL, nil
			}
		}
		// If no PDF type found, return the first item
		if signedURLResp.Items[0].SignedURL != "" {
			return signedURLResp.Items[0].SignedURL, nil
		}
	}

	return "", fmt.Errorf("no download URL found in response")
}

// DownloadFile downloads a file from the given URL
func (c *Client) DownloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// FindBillingByInvoiceID searches for a billing that contains the specified invoice ID
// It searches through recent billings (last 2 years) to find the matching invoice
func (c *Client) FindBillingByInvoiceID(invoiceID string) (*Billing, error) {
	// Search in the last 2 years to find the billing
	// This is a reasonable time range for invoice lookup
	now := time.Now()
	twoYearsAgo := now.AddDate(-2, 0, 0)

	startDate := twoYearsAgo.Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	billings, err := c.GetBillings(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get billings: %w", err)
	}

	// Search for the invoice ID in all billings
	for i := range billings {
		for _, invID := range billings[i].InvoiceIDs {
			if invID == invoiceID {
				return &billings[i], nil
			}
		}
	}

	return nil, fmt.Errorf("invoice ID %s not found in any billing", invoiceID)
}
