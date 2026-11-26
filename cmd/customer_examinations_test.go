package cmd

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestCustomerExaminationsCreateHasAllPayloadFlags(t *testing.T) {
	// POST /customer_examinations の項目（swagger の CustomerExaminationPayload）にすべてフラグがあること
	for _, name := range []string{
		"customer-id", "amount", "end-date", "address1", "address2", "email", "tel", "zip-code",
		"business-type", "corporate-number", "representative-name", "business-description",
		"website", "remark", "from-destination",
	} {
		if customerExaminationsCreateCmd.Flags().Lookup(name) == nil {
			t.Errorf("--%s が無い", name)
		}
	}
}

func TestFillFromDestinationKeepsExplicitValues(t *testing.T) {
	body := map[string]any{"amount": 3000000, "tel": "03-0000-0000"}
	dest := map[string]any{
		"id": "DEST", "customer_id": "CUST", "address1": "東京都千代田区", "address2": "",
		"zip_code": "100-0001", "tel": "06-1111-1111", "email": "billing@example.com", "name": "担当者",
	}
	fillFromDestination(body, dest)
	want := map[string]any{
		"amount": 3000000, "tel": "03-0000-0000", "customer_id": "CUST",
		"address1": "東京都千代田区", "zip_code": "100-0001", "email": "billing@example.com",
	}
	if len(body) != len(want) {
		t.Fatalf("body = %v, want %v", body, want)
	}
	for k, v := range want {
		if body[k] != v {
			t.Errorf("%s = %v, want %v", k, body[k], v)
		}
	}
}

func TestCustomerExaminationsCreateDryRunBody(t *testing.T) {
	t.Setenv("MFK_API_KEY", "dummy")
	t.Setenv("MFK_ENV", "")
	orig := destinationFetcher
	destinationFetcher = func(id string) (map[string]any, error) {
		if id != "DEST" {
			t.Errorf("destination id = %q", id)
		}
		return map[string]any{
			"customer_id": "CUST", "address1": "東京都千代田区", "address2": "1F",
			"zip_code": "100-0001", "tel": "03-0000-0000", "email": "billing@example.com",
		}, nil
	}
	defer func() { destinationFetcher = orig; flagDryRun = false }()

	r, w, _ := os.Pipe()
	stderr := os.Stderr
	os.Stderr = w
	rootCmd.SetArgs([]string{
		"customer-examinations", "create", "--dry-run",
		"--from-destination", "DEST", "--amount", "3000000", "--end-date", "2026-12-31",
		"--email", "override@example.com", "--business-type", "corporate",
		"--corporate-number", "1234567890123", "--website", "https://example.com",
		"--business-description", "Web services", "--remark", "note",
	})
	err := rootCmd.Execute()
	w.Close()
	os.Stderr = stderr
	rootCmd.SetArgs(nil)
	if err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(r)

	var req struct {
		Method string         `json:"method"`
		URL    string         `json:"url"`
		Body   map[string]any `json:"body"`
	}
	if err := json.Unmarshal(out, &req); err != nil {
		t.Fatalf("dry-run の出力が JSON でない: %v\n%s", err, out)
	}
	if req.Method != "POST" || req.URL != "https://api.mfkessai.co.jp/v2/customer_examinations" {
		t.Errorf("request = %s %s", req.Method, req.URL)
	}
	want := map[string]any{
		"customer_id": "CUST", "amount": float64(3000000), "end_date": "2026-12-31",
		"address1": "東京都千代田区", "address2": "1F", "zip_code": "100-0001", "tel": "03-0000-0000",
		"email": "override@example.com", "business_type": "corporate", "corporate_number": "1234567890123",
		"website": "https://example.com", "business_description": "Web services", "remark": "note",
	}
	if len(req.Body) != len(want) {
		t.Errorf("body = %v", req.Body)
	}
	for k, v := range want {
		if req.Body[k] != v {
			t.Errorf("%s = %v, want %v", k, req.Body[k], v)
		}
	}
}
