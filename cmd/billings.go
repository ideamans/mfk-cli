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

var billingsUploadSignedURLCmd = &cobra.Command{
	Use:   "upload-signed-url",
	Short: "Get upload signed URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := buildBody(nil)
		if err != nil {
			return err
		}
		return doPost("/billings/upload_signed_url", body)
	},
}

var billingsDownloadSignedURLCmd = &cobra.Command{
	Use:   "download-signed-url",
	Short: "Get download signed URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := buildBody(nil)
		if err != nil {
			return err
		}
		return doPost("/billings/download_signed_url", body)
	},
}
