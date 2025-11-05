package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/miyanaga/moneyforward-kessai-invoice-downloader-v2/internal/downloader"
	"github.com/miyanaga/moneyforward-kessai-invoice-downloader-v2/internal/mfkessai"
	"github.com/spf13/cobra"
)

var (
	concurrency int
)

var downloadInvoicesCmd = &cobra.Command{
	Use:   "download-invoices [開始日] [終了日] [出力ディレクトリ]",
	Short: "指定期間の請求書PDFをダウンロード",
	Long: `指定した期間内の請求書PDFファイルをダウンロードします。
日付はYYYY-MM-DD形式で指定してください。
出力ディレクトリが存在しない場合は自動的に作成されます。

使用例:
  mfk -c 5 download-invoices 2024-09-01 2025-09-30 ~/Downloads/mfk`,
	Args: cobra.ExactArgs(3),
	RunE: runDownloadInvoices,
}

func init() {
	downloadInvoicesCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 5, "並列ダウンロード数")
}

func runDownloadInvoices(cmd *cobra.Command, args []string) error {
	startDate := args[0]
	endDate := args[1]
	outputDir := args[2]

	// Determine API key: command-line flag > environment variable
	var finalAPIKey string
	if apiKey != "" {
		// Use API key from command-line flag
		finalAPIKey = apiKey
	} else {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}

		// Get API key from environment
		finalAPIKey = os.Getenv("MFK_API_KEY")
		if finalAPIKey == "" {
			return fmt.Errorf("API key not provided: use --api-key flag or set MFK_API_KEY environment variable")
		}
	}

	// Resolve tilde in output directory
	if strings.HasPrefix(outputDir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		outputDir = filepath.Join(home, outputDir[2:])
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	log.Printf("Starting invoice download")
	log.Printf("Date range: %s to %s", startDate, endDate)
	log.Printf("Output directory: %s", outputDir)
	log.Printf("Concurrency: %d", concurrency)

	// Create API client
	client := mfkessai.NewClient(finalAPIKey)

	// Create downloader
	dl := downloader.New(client, concurrency)

	// Download invoices
	if err := dl.DownloadInvoices(startDate, endDate, outputDir); err != nil {
		return fmt.Errorf("failed to download invoices: %w", err)
	}

	log.Println("Download completed successfully")
	return nil
}
