package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(customersCmd)
	customersCmd.AddCommand(customersListCmd, customersGetCmd, customersCreateCmd, customersUpdateCmd, customersBankTransferCmd, customersAccountTransferCmd)

	customersListCmd.Flags().String("number", "", "Filter by customer number")
	customersListCmd.Flags().String("name", "", "Filter by customer name")
	customersListCmd.Flags().Bool("has-alert", false, "Filter by alert status")

	customersCreateCmd.Flags().String("name", "", "Customer name")
	customersCreateCmd.Flags().String("number", "", "Customer number")
	customersCreateCmd.Flags().String("dest-name", "", "Destination name")
	customersCreateCmd.Flags().String("dest-name-kana", "", "Destination name kana")
	customersCreateCmd.Flags().String("dest-email", "", "Destination email")
	customersCreateCmd.Flags().String("dest-tel", "", "Destination phone")
	customersCreateCmd.Flags().String("dest-zip-code", "", "Destination zip code")
	customersCreateCmd.Flags().String("dest-address1", "", "Destination address1")
	customersCreateCmd.Flags().String("dest-address2", "", "Destination address2")
	customersCreateCmd.Flags().String("dest-department", "", "Destination department")
	customersCreateCmd.Flags().String("dest-title", "", "Destination title")

	customersUpdateCmd.Flags().String("name", "", "Customer name")
}

var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Manage customers",
}

var customersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List customers",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("number"); v != "" {
			query.Set("customer_number", v)
		}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			query.Set("customer_name", v)
		}
		if v, _ := cmd.Flags().GetBool("has-alert"); v {
			query.Set("has_alert", "true")
		}
		return doList("/customers", query)
	},
}

var customersGetCmd = &cobra.Command{
	Use:   "get <customer_id>",
	Short: "Get customer details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/customers/%s", args[0]))
	},
}

var customersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a customer",
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			overrides["name"] = v
		}
		if v, _ := cmd.Flags().GetString("number"); v != "" {
			overrides["number"] = v
		}
		dest := map[string]any{}
		destFields := map[string]string{
			"dest-name": "name", "dest-name-kana": "name_kana", "dest-email": "email",
			"dest-tel": "tel", "dest-zip-code": "zip_code", "dest-address1": "address1",
			"dest-address2": "address2", "dest-department": "department", "dest-title": "title",
		}
		for flag, field := range destFields {
			if v, _ := cmd.Flags().GetString(flag); v != "" {
				dest[field] = v
			}
		}
		if len(dest) > 0 {
			overrides["destination"] = dest
		}
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost("/customers", body)
	},
}

var customersUpdateCmd = &cobra.Command{
	Use:   "update <customer_id>",
	Short: "Update a customer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			overrides["name"] = v
		}
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPatch(fmt.Sprintf("/customers/%s", args[0]), body)
	},
}

var customersBankTransferCmd = &cobra.Command{
	Use:   "bank-transfer <customer_id>",
	Short: "Assign bank transfer account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := buildBody(nil)
		if err != nil {
			return err
		}
		return doPost(fmt.Sprintf("/customers/%s/bank_transfers", args[0]), body)
	},
}

var customersAccountTransferCmd = &cobra.Command{
	Use:   "account-transfer <customer_id>",
	Short: "Apply for account transfer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := buildBody(nil)
		if err != nil {
			return err
		}
		return doPost(fmt.Sprintf("/customers/%s/account_transfer_requests", args[0]), body)
	},
}
