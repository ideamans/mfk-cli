package cmd

import "testing"

func TestBillingDownloadSignedURLPath(t *testing.T) {
	got := billingDownloadSignedURLPath("BILLING", "INVOICE")
	want := "/billings/BILLING/issues/i/INVOICE/download_signed_url"
	if got != want {
		t.Errorf("path = %q, want %q", got, want)
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
