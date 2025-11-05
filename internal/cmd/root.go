package cmd

import (
	"github.com/spf13/cobra"
)

var (
	apiKey string
)

var rootCmd = &cobra.Command{
	Use:   "mfk",
	Short: "マネーフォワードケッサイ CLIツール",
	Long:  "マネーフォワードケッサイAPIを操作するためのコマンドラインツール",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "マネーフォワードケッサイAPIキー（環境変数MFK_API_KEYより優先）")
	rootCmd.AddCommand(downloadInvoiceCmd)
	rootCmd.AddCommand(downloadInvoicesCmd)
}
