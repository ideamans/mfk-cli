package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(customerNameUpdatesCmd)
	customerNameUpdatesCmd.AddCommand(customerNameUpdatesListCmd, customerNameUpdatesGetCmd, customerNameUpdatesCreateCmd, customerNameUpdatesDeleteCmd)

	customerNameUpdatesCreateCmd.Flags().String("customer-id", "", "Customer ID")
	customerNameUpdatesCreateCmd.Flags().String("name", "", "New customer name")
}

var customerNameUpdatesCmd = &cobra.Command{
	Use:   "customer-name-updates",
	Short: "Manage customer name updates",
}

var customerNameUpdatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List customer name updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doList("/customer_name_updates", url.Values{})
	},
}

var customerNameUpdatesGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get customer name update details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/customer_name_updates/%s", args[0]))
	},
}

var customerNameUpdatesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a customer name update",
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		setIfNotEmpty(overrides, "customer_id", mustGetString(cmd, "customer-id"))
		setIfNotEmpty(overrides, "name", mustGetString(cmd, "name"))
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost("/customer_name_updates", body)
	},
}

var customerNameUpdatesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Cancel a customer name update",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doDelete(fmt.Sprintf("/customer_name_updates/%s", args[0]))
	},
}
