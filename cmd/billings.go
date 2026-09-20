package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(billingsCmd)
	billingsCmd.AddCommand(billingsListCmd, billingsQualifiedCmd, billingsGetCmd, billingsReissueCmd, billingsUploadSignedURLCmd, billingsDownloadSignedURLCmd)

	billingsListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	billingsListCmd.Flags().String("destination-id", "", "Filter by destination ID")
	billingsListCmd.Flags().String("status", "", "Filter by status")

	billingsQualifiedCmd.Flags().String("customer-id", "", "Filter by customer ID")
	billingsQualifiedCmd.Flags().String("destination-id", "", "Filter by destination ID")

	billingsUploadSignedURLCmd.Flags().String("content-type", "application/pdf",
		"Content type of the file to upload: application/pdf, application/json, text/csv or text/plain")
}

var billingsCmd = &cobra.Command{
	Use:   "billings",
	Short: "Manage billings (use `qualified` for invoice-system billings since Oct 2023)",
	Long: `Manage billings.

Money Forward Kessai splits billing retrieval into two endpoints around the
Japanese qualified-invoice system (インボイス制度, started Oct 2023):

  mfk billings qualified  -> GET /billings/qualified
      Qualified invoice system (適格請求書等保存方式). Billings from Oct 2023
      onward. THIS IS THE MAIN COMMAND for current billings.

  mfk billings list       -> GET /billings
      Legacy classification-style billings (区分記載請求書等保存方式), roughly
      before Oct 2023. Use this only to access older billings.

Add --page-all to fetch every page (otherwise only the latest 20 are returned).`,
}

var billingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List legacy (pre-Oct-2023) billings — use `qualified` for current invoice-system billings",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("customer-id"); v != "" {
			query.Set("customer_id", v)
		}
		if v, _ := cmd.Flags().GetString("destination-id"); v != "" {
			query.Set("destination_id", v)
		}
		if v, _ := cmd.Flags().GetString("status"); v != "" {
			query.Set("status", v)
		}
		return doList("/billings", query)
	},
}

var billingsQualifiedCmd = &cobra.Command{
	Use:   "qualified",
	Short: "List invoice-system billings (適格請求書, since Oct 2023) — the main billings command",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("customer-id"); v != "" {
			query.Set("customer_id", v)
		}
		if v, _ := cmd.Flags().GetString("destination-id"); v != "" {
			query.Set("destination_id", v)
		}
		return doList("/billings/qualified", query)
	},
}

var billingsGetCmd = &cobra.Command{
	Use:   "get <billing_id>",
	Short: "Get billing details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/billings/%s", args[0]))
	},
}

var billingsReissueCmd = &cobra.Command{
	Use:   "reissue <billing_id>",
	Short: "Reissue invoice",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := buildBody(nil)
		if err != nil {
			return err
		}
		return doPost(fmt.Sprintf("/billings/%s/reissue", args[0]), body)
	},
}

// billingUploadSignedURLPath はファイルアップロード用署名URLを発行するパス。
// 請求（billing）ごとに発行するので請求IDが要る。
// 以前は /billings/upload_signed_url に投げていたが、そのルートは存在せず 405 になる。
func billingUploadSignedURLPath(billingID string) string {
	return fmt.Sprintf("/billings/%s/upload_signed_url", billingID)
}

var billingsUploadSignedURLCmd = &cobra.Command{
	Use:   "upload-signed-url <billing_id>",
	Short: "Get a signed URL to upload a file for one billing",
	Long: `Get a short-lived signed URL for uploading a file attached to one billing.

--content-type is required by the API and must be one of application/pdf,
application/json, text/csv or text/plain (default: application/pdf).

This endpoint is only available to sellers whose contract includes it; other
accounts get "forbidden_seller" (HTTP 403) even though the request is correct.
Ask Money Forward Kessai to enable it before using this command.`,
	Example: `  mfk billings upload-signed-url <billing_id>
  mfk billings upload-signed-url <billing_id> --content-type text/csv`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		contentType, _ := cmd.Flags().GetString("content-type")
		setIfNotEmpty(overrides, "content_type", contentType)
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost(billingUploadSignedURLPath(args[0]), body)
	},
}

// billingDownloadSignedURLPath は請求書PDFのダウンロード用署名URLを発行するパス。
// 請求（billing）に請求書（invoice）が紐づくので、両方のIDが要る。
// 以前は /billings/download_signed_url に投げていたが、そのルートは存在せず 405 になる。
func billingDownloadSignedURLPath(billingID, invoiceID string) string {
	return fmt.Sprintf("/billings/%s/issues/i/%s/download_signed_url", billingID, invoiceID)
}

var billingsDownloadSignedURLCmd = &cobra.Command{
	Use:   "download-signed-url <billing_id> <invoice_id>",
	Short: "Get a signed URL to download one invoice PDF",
	Long: `Get a short-lived signed URL for downloading the PDF of one invoice.

The invoice ID is in the billing's invoice_ids (see "mfk billings get" or
"mfk billings qualified"). The response looks like
{"items":[{"signed_url": "...", "expired_at": "...", "type": "pdf"}]}.
Download the file from signed_url without the API key; the URL expires soon.`,
	Example: `  mfk billings download-signed-url <billing_id> <invoice_id>
  curl -sSo invoice.pdf "$(mfk billings download-signed-url <billing_id> <invoice_id> | jq -r '.items[0].signed_url')"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doPost(billingDownloadSignedURLPath(args[0], args[1]), nil)
	},
}
