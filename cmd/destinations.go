package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var destFlags = []struct {
	flag, field string
}{
	{"name", "name"}, {"name-kana", "name_kana"}, {"email", "email"},
	{"tel", "tel"}, {"zip-code", "zip_code"}, {"address1", "address1"},
	{"address2", "address2"}, {"department", "department"}, {"title", "title"},
}

func init() {
	rootCmd.AddCommand(destinationsCmd)
	destinationsCmd.AddCommand(destinationsListCmd, destinationsGetCmd, destinationsCreateCmd, destinationsUpdateCmd)

	destinationsListCmd.Flags().String("customer-id", "", "Filter by customer ID")

	destinationsCreateCmd.Flags().String("customer-id", "", "Customer ID")
	for _, f := range destFlags {
		destinationsCreateCmd.Flags().String(f.flag, "", "Destination "+f.field)
	}

	for _, f := range destFlags {
		destinationsUpdateCmd.Flags().String(f.flag, "", "Destination "+f.field)
	}
}

var destinationsCmd = &cobra.Command{
	Use:   "destinations",
	Short: "Manage billing destinations",
}

var destinationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List destinations",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("customer-id"); v != "" {
			query.Set("customer_id", v)
		}
		return doList("/destinations", query)
	},
}

var destinationsGetCmd = &cobra.Command{
	Use:   "get <destination_id>",
	Short: "Get destination details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doGet(fmt.Sprintf("/destinations/%s", args[0]))
	},
}

var destinationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a destination",
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		if v, _ := cmd.Flags().GetString("customer-id"); v != "" {
			overrides["customer_id"] = v
		}
		for _, f := range destFlags {
			if v, _ := cmd.Flags().GetString(f.flag); v != "" {
				overrides[f.field] = v
			}
		}
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPost("/destinations", body)
	},
}

var destinationsUpdateCmd = &cobra.Command{
	Use:   "update <destination_id>",
	Short: "Update a destination",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overrides := map[string]any{}
		for _, f := range destFlags {
			if v, _ := cmd.Flags().GetString(f.flag); v != "" {
				overrides[f.field] = v
			}
		}
		body, err := buildBody(overrides)
		if err != nil {
			return err
		}
		return doPatch(fmt.Sprintf("/destinations/%s", args[0]), body)
	},
}
