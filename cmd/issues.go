package cmd

import (
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(issuesCmd)
	issuesCmd.AddCommand(issuesListCmd)
}

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage issued billing information",
}

var issuesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issued billing information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doList("/issues", url.Values{})
	},
}
