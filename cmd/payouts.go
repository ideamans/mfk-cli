package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(payoutsCmd)
	payoutsCmd.AddCommand(payoutsListCmd, payoutsGetCmd)
}

var payoutsCmd = &cobra.Command{
	Use:   "payouts",
	Short: "Manage payouts",
}

var payoutsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List payouts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doList("/payouts", url.Values{})
	},
}

var payoutsGetCmd = &cobra.Command{
	Use:   "get <payout_id>",
	Short: "Get payout details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/payouts/%s", args[0]))
	},
}
