package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(payoutTransactionsCmd)
	payoutTransactionsCmd.AddCommand(payoutTransactionsListCmd, payoutTransactionsGetCmd)

	payoutTransactionsListCmd.Flags().String("payout-id", "", "Filter by payout ID")
}

var payoutTransactionsCmd = &cobra.Command{
	Use:   "payout-transactions",
	Short: "Manage payout transactions",
}

var payoutTransactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List payout transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("payout-id"); v != "" {
			query.Set("payout_id", v)
		}
		return doList("/payout_transactions", query)
	},
}

var payoutTransactionsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get payout transaction details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/payout_transactions/%s", args[0]))
	},
}
