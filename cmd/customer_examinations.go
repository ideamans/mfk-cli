package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(customerExaminationsCmd)
	customerExaminationsCmd.AddCommand(customerExaminationsListCmd, customerExaminationsGetCmd, customerExaminationsCreateCmd)

	customerExaminationsListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	customerExaminationsListCmd.Flags().String("status", "", "Filter by status")

	customerExaminationsCreateCmd.Flags().String("customer-id", "", "Customer ID")
	customerExaminationsCreateCmd.Flags().Int("amount", 0, "Amount")
	customerExaminationsCreateCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	customerExaminationsCreateCmd.Flags().String("address1", "", "Address")
	customerExaminationsCreateCmd.Flags().String("email", "", "Email")
	customerExaminationsCreateCmd.Flags().String("tel", "", "Phone")
	customerExaminationsCreateCmd.Flags().String("zip-code", "", "Zip code")
	customerExaminationsCreateCmd.Flags().String("business-type", "", "Business type")
	customerExaminationsCreateCmd.Flags().String("representative-name", "", "Representative name")
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

var customerExaminationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a customer examination",
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		setIfNotEmpty(overrides, "customer_id", mustGetString(cmd, "customer-id"))
		setIfNotZero(overrides, "amount", mustGetInt(cmd, "amount"))
		setIfNotEmpty(overrides, "end_date", mustGetString(cmd, "end-date"))
		setIfNotEmpty(overrides, "address1", mustGetString(cmd, "address1"))
		setIfNotEmpty(overrides, "email", mustGetString(cmd, "email"))
		setIfNotEmpty(overrides, "tel", mustGetString(cmd, "tel"))
		setIfNotEmpty(overrides, "zip_code", mustGetString(cmd, "zip-code"))
		setIfNotEmpty(overrides, "business_type", mustGetString(cmd, "business-type"))
		setIfNotEmpty(overrides, "representative_name", mustGetString(cmd, "representative-name"))
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost("/customer_examinations", body)
	},
}

func mustGetString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func mustGetInt(cmd *cobra.Command, name string) int {
	v, _ := cmd.Flags().GetInt(name)
	return v
}
