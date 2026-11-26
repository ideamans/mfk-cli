package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(customerExaminationsCmd)
	customerExaminationsCmd.AddCommand(customerExaminationsListCmd, customerExaminationsGetCmd, customerExaminationsCreateCmd)

	customerExaminationsListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	customerExaminationsListCmd.Flags().String("status", "", "Filter by status")

	f := customerExaminationsCreateCmd.Flags()
	f.String("customer-id", "", "Customer ID (required; filled from --from-destination if omitted)")
	f.Int("amount", 0, "Requested credit amount in JPY (required). Must not be below the credit line currently in effect")
	f.String("end-date", "", "Last day transactions may be registered against this credit line (required, YYYY-MM-DD). Only the last day of one of the next three months is accepted")
	f.String("address1", "", "Address line 1 (required)")
	f.String("address2", "", "Address line 2 (building, floor)")
	f.String("email", "", "Email (required)")
	f.String("tel", "", "Phone (required)")
	f.String("zip-code", "", "Zip code (required)")
	f.String("business-type", "", "Business type: corporate or individual. Leave empty if unknown")
	f.String("corporate-number", "", "Corporate number (法人番号). Only when --business-type corporate")
	f.String("representative-name", "", "Representative name")
	f.String("business-description", "", "Business description: the customer's main services or products")
	f.String("website", "", "Customer's website URL")
	f.String("remark", "", "Other information for the examination. In sandbox, also selects the examination result")
	f.String("from-destination", "", "Fill customer_id, address1, address2, zip_code, tel and email from this billing destination (GET /destinations/{id}). Explicit flags and --json win")
}

var customerExaminationsCmd = &cobra.Command{
	Use:   "customer-examinations",
	Short: "Manage customer examinations",
}

var customerExaminationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List customer examinations",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("customer-id"); v != "" {
			query.Set("customer_id", v)
		}
		if v, _ := cmd.Flags().GetString("status"); v != "" {
			query.Set("status", v)
		}
		return doList("/customer_examinations", query)
	},
}

var customerExaminationsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get customer examination details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/customer_examinations/%s", args[0]))
	},
}

// customerExaminationFlags は create のフラグとリクエストの項目の対応です。
var customerExaminationFlags = []struct {
	flag, field string
}{
	{"customer-id", "customer_id"}, {"end-date", "end_date"},
	{"address1", "address1"}, {"address2", "address2"},
	{"email", "email"}, {"tel", "tel"}, {"zip-code", "zip_code"},
	{"business-type", "business_type"}, {"corporate-number", "corporate_number"},
	{"representative-name", "representative_name"},
	{"business-description", "business_description"},
	{"website", "website"}, {"remark", "remark"},
}

// destinationExaminationFields は --from-destination で請求先から写す項目です。
var destinationExaminationFields = []string{"customer_id", "address1", "address2", "zip_code", "tel", "email"}

var customerExaminationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Request a credit examination (also used to request a higher credit line)",
	Long: `Request a credit examination (与信枠審査) for a customer: POST /customer_examinations.

This is also how you request a higher credit line (増枠申請). The examination
finishes within two business days at most.

Required: customer_id, amount, end_date, address1, email, tel, zip_code.

  --amount     The requested credit line in JPY. While a credit line is in
               effect, an amount below it is rejected. The examination may
               grant less than requested.
  --end-date   Transactions can be registered against the line until this
               date. Only the last day of one of the next three months is
               accepted (e.g. 2026-12-31).

With automatic credit examination (自動与信), the raised amount carries over
to later automatic credit lines, and the customer's alert is cleared.

While an examination for the customer is still in progress, a new one is
rejected with HTTP 409.

--from-destination <destination_id> fetches an existing billing destination
(a read-only GET, sent even with --dry-run) and fills customer_id, address1,
address2, zip_code, tel and email from it. Explicit flags win over --json,
and --json wins over the destination. A credit line increase then needs
little more than:

  mfk customer-examinations create --from-destination <destination_id> \
    --amount 3000000 --end-date 2026-12-31 --dry-run

Always check the request with --dry-run first: this sends a real request
to Money Forward Kessai.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		for _, f := range customerExaminationFlags {
			setIfNotEmpty(overrides, f.field, mustGetString(cmd, f.flag))
		}
		setIfNotZero(overrides, "amount", mustGetInt(cmd, "amount"))
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		if destID := mustGetString(cmd, "from-destination"); destID != "" {
			dest, err := destinationFetcher(destID)
			if err != nil {
				return err
			}
			fillFromDestination(body, dest)
		}
		return doPost("/customer_examinations", body)
	},
}

// destinationFetcher はテストで差し替えます。
var destinationFetcher = fetchDestination

// fetchDestination は請求先を取得します。読み取りなので --dry-run でも送ります。
func fetchDestination(id string) (map[string]any, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	client.DryRun = false
	data, err := client.Get(fmt.Sprintf("/destinations/%s", url.PathEscape(id)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch destination %s: %w", id, err)
	}
	var dest map[string]any
	if err := json.Unmarshal(data, &dest); err != nil {
		return nil, fmt.Errorf("failed to parse destination %s: %w", id, err)
	}
	return dest, nil
}

// fillFromDestination は body にまだ無い項目だけを請求先の値で埋めます。
func fillFromDestination(body, dest map[string]any) {
	for _, k := range destinationExaminationFields {
		if _, ok := body[k]; ok {
			continue
		}
		if v, ok := dest[k].(string); ok && v != "" {
			body[k] = v
		}
	}
}

func mustGetString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func mustGetInt(cmd *cobra.Command, name string) int {
	v, _ := cmd.Flags().GetInt(name)
	return v
}
