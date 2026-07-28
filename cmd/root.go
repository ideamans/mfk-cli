package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ideamans/go-llm-cli-kit/llmcmd"
	"github.com/spf13/cobra"

	"github.com/ideamans/mfk-cli/pkg/api"
	"github.com/ideamans/mfk-cli/pkg/output"
)

var (
	flagFormat  string
	flagSandbox bool
	flagPageAll bool
	flagLimit   int
	flagDryRun  bool
	flagJSON    string
)

// PluginVersion はこのCLIのリリースバージョンです。
// plugins/mfk-cli/.claude-plugin/plugin.json の version と一致していることを
// テストが、git タグと一致していることをリリースワークフローが検査します。
const PluginVersion = "0.4.0"

var rootCmd = &cobra.Command{
	Use:     "mfk",
	Short:   "Money Forward Kessai CLI",
	Long:    "CLI tool for Money Forward Kessai API v2",
	Version: PluginVersion,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Root は組み立て済みのコマンドツリーを実行せずに返します。
// カタログ生成器が Execute と同じ定義から生成するために使います。
func Root() *cobra.Command { return rootCmd }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagFormat, "format", "json", "Output format: json, table, csv")
	rootCmd.PersistentFlags().BoolVar(&flagSandbox, "sandbox", false, "Use sandbox environment")
	rootCmd.PersistentFlags().BoolVar(&flagPageAll, "page-all", false, "Fetch all pages")
	rootCmd.PersistentFlags().IntVar(&flagLimit, "limit", 0, "Number of items to fetch")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "Show request without executing")
	rootCmd.PersistentFlags().StringVar(&flagJSON, "json", "", "Request body as JSON string")
	llmcmd.AddTo(rootCmd, LLMConfig())
}

func newClient() (*api.Client, error) {
	sandbox, err := api.IsSandbox(flagSandbox)
	if err != nil {
		return nil, err
	}
	client, err := api.NewClient(sandbox)
	if err != nil {
		return nil, err
	}
	client.DryRun = flagDryRun
	return client, nil
}

func printResult(data []byte) error {
	return output.Print(data, flagFormat)
}

func applyPagination(query url.Values) url.Values {
	if query == nil {
		query = url.Values{}
	}
	if flagLimit > 0 {
		query.Set("limit", fmt.Sprintf("%d", flagLimit))
	}
	return query
}

func buildBody(overrides map[string]any) (map[string]any, error) {
	return api.BuildBody(flagJSON, overrides)
}

func setIfNotEmpty(m map[string]any, key, value string) {
	if value != "" {
		m[key] = value
	}
}

func setIfNotZero(m map[string]any, key string, value int) {
	if value != 0 {
		m[key] = value
	}
}

func doList(path string, query url.Values) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	query = applyPagination(query)

	var data []byte
	if flagPageAll {
		data, err = client.GetAllPages(path, query)
	} else {
		data, err = client.Get(path, query)
	}
	if err != nil {
		return err
	}
	return printResult(data)
}

func doGet(path string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	data, err := client.Get(path, nil)
	if err != nil {
		return err
	}
	return printResult(data)
}

func doPost(path string, body map[string]any) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	data, err := client.Post(path, body)
	if err != nil {
		return err
	}
	return printResult(data)
}

func doPatch(path string, body map[string]any) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	data, err := client.Patch(path, body)
	if err != nil {
		return err
	}
	return printResult(data)
}

func doDelete(path string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	data, err := client.Delete(path)
	if err != nil {
		return err
	}
	return printResult(data)
}
