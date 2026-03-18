package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(creditFacilitiesCmd)
	creditFacilitiesCmd.AddCommand(creditFacilitiesListCmd, creditFacilitiesGetCmd)

	creditFacilitiesListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	creditFacilitiesListCmd.Flags().String("customer-examination-id", "", "Filter by examination ID")
}

var creditFacilitiesCmd = &cobra.Command{
	Use:   "credit-facilities",
	Short: "Manage credit facilities",
}

var creditFacilitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List credit facilities",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("customer-id"); v != "" {
			query.Set("customer_id", v)
		}
		if v, _ := cmd.Flags().GetString("customer-examination-id"); v != "" {
			query.Set("customer_examination_id", v)
		}
		return doList("/credit_facilities", query)
	},
}

var creditFacilitiesGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get credit facility details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/credit_facilities/%s", args[0]))
	},
}
