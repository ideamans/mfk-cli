package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authorizationsCmd)
	authorizationsCmd.AddCommand(authorizationsListCmd, authorizationsGetCmd, authorizationsCreateCmd, authorizationsDeleteCmd)

	authorizationsCreateCmd.Flags().String("customer-id", "", "Customer ID")
	authorizationsCreateCmd.Flags().Int("amount", 0, "Amount")
	authorizationsCreateCmd.Flags().String("destination-id", "", "Destination ID")
}

var authorizationsCmd = &cobra.Command{
	Use:   "authorizations",
	Short: "Manage authorizations",
}

var authorizationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List authorizations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doList("/authorizations", url.Values{})
	},
}

var authorizationsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get authorization details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/authorizations/%s", args[0]))
	},
}

var authorizationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an authorization",
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		setIfNotEmpty(overrides, "customer_id", mustGetString(cmd, "customer-id"))
		setIfNotZero(overrides, "amount", mustGetInt(cmd, "amount"))
		setIfNotEmpty(overrides, "destination_id", mustGetString(cmd, "destination-id"))
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost("/authorizations", body)
	},
}

var authorizationsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Cancel an authorization",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doDelete(fmt.Sprintf("/authorizations/%s", args[0]))
	},
}
