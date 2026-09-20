package cmd

import "testing"

func TestBillingDownloadSignedURLPath(t *testing.T) {
	got := billingDownloadSignedURLPath("BILLING", "INVOICE")
	want := "/billings/BILLING/issues/i/INVOICE/download_signed_url"
	if got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
}

func TestBillingUploadSignedURLPath(t *testing.T) {
	got := billingUploadSignedURLPath("BILLING")
	want := "/billings/BILLING/upload_signed_url"
	if got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
}

func TestBillingsUploadSignedURLNeedsBillingID(t *testing.T) {
	cmd := billingsUploadSignedURLCmd
	if err := cmd.Args(cmd, nil); err == nil {
		t.Error("請求IDなしで通ってしまう。請求IDは必須のはず")
	}
	if err := cmd.Args(cmd, []string{"BILLING"}); err != nil {
		t.Errorf("請求IDで通らない: %v", err)
	}
	if got, _ := cmd.Flags().GetString("content-type"); got != "application/pdf" {
		t.Errorf("content-type の既定値 = %q, want application/pdf", got)
	}
}

func TestBillingsDownloadSignedURLNeedsBothIDs(t *testing.T) {
	cmd := billingsDownloadSignedURLCmd
	if err := cmd.Args(cmd, []string{"BILLING"}); err == nil {
		t.Error("請求IDだけで通ってしまう。請求書IDも必須のはず")
	}
	if err := cmd.Args(cmd, []string{"BILLING", "INVOICE"}); err != nil {
		t.Errorf("請求IDと請求書IDで通らない: %v", err)
	}
}
