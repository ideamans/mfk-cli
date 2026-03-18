package cmd

import (
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(payoutRefundsCmd)
	payoutRefundsCmd.AddCommand(payoutRefundsListCmd)
}

var payoutRefundsCmd = &cobra.Command{
	Use:   "payout-refunds",
	Short: "Manage payout refunds",
}

var payoutRefundsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List payout refunds",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doList("/payout_refunds", url.Values{})
	},
}
