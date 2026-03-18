package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(transactionsCmd)
	transactionsCmd.AddCommand(transactionsListCmd, transactionsGetCmd, transactionsCreateCmd, transactionsDeleteCmd)

	transactionsListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	transactionsListCmd.Flags().String("destination-id", "", "Filter by destination ID")
	transactionsListCmd.Flags().String("status", "", "Filter by status")

	transactionsCreateCmd.Flags().String("destination-id", "", "Destination ID")
	transactionsCreateCmd.Flags().Int("amount", 0, "Transaction amount")
	transactionsCreateCmd.Flags().String("number", "", "Transaction number")
	transactionsCreateCmd.Flags().String("date", "", "Transaction date (YYYY-MM-DD)")
	transactionsCreateCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD)")
	transactionsCreateCmd.Flags().String("issue-date", "", "Issue date (YYYY-MM-DD)")
	transactionsCreateCmd.Flags().String("delivery-method", "", "Invoice delivery method: posting or email")
}

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Manage transactions",
}

var transactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions",
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
		return doList("/transactions", query)
	},
}

var transactionsGetCmd = &cobra.Command{
	Use:   "get <transaction_id>",
	Short: "Get transaction details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/transactions/%s", args[0]))
	},
}

var transactionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		if v, _ := cmd.Flags().GetString("destination-id"); v != "" {
			overrides["destination_id"] = v
		}
		if v, _ := cmd.Flags().GetInt("amount"); v != 0 {
			overrides["amount"] = v
		}
		if v, _ := cmd.Flags().GetString("number"); v != "" {
			overrides["number"] = v
		}
		if v, _ := cmd.Flags().GetString("date"); v != "" {
			overrides["date"] = v
		}
		if v, _ := cmd.Flags().GetString("due-date"); v != "" {
			overrides["due_date"] = v
		}
		if v, _ := cmd.Flags().GetString("issue-date"); v != "" {
			overrides["issue_date"] = v
		}
		if v, _ := cmd.Flags().GetString("delivery-method"); v != "" {
			overrides["invoice_delivery_methods"] = []string{v}
		}
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost("/transactions", body)
	},
}

var transactionsDeleteCmd = &cobra.Command{
	Use:   "delete <transaction_id>",
	Short: "Cancel a transaction",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doDelete(fmt.Sprintf("/transactions/%s", args[0]))
	},
}
