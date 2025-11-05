package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/miyanaga/moneyforward-kessai-invoice-downloader-v2/internal/mfkessai"
	"github.com/spf13/cobra"
)

var downloadInvoiceCmd = &cobra.Command{
	Use:   "download-invoice [請求書ID] [出力ディレクトリ]",
	Short: "指定した請求書IDのPDFをダウンロード",
	Long: `指定した請求書IDのPDFファイルをダウンロードします。
出力ディレクトリが存在しない場合は自動的に作成されます。

使用例:
  mfk download-invoice IN_XXXXXXXXXXXXX ~/Downloads/mfk`,
	Args: cobra.ExactArgs(2),
	RunE: runDownloadInvoice,
}

func runDownloadInvoice(cmd *cobra.Command, args []string) error {
	invoiceID := args[0]
	outputDir := args[1]

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
	log.Printf("Invoice ID: %s", invoiceID)
	log.Printf("Output directory: %s", outputDir)

	// Create API client
	client := mfkessai.NewClient(finalAPIKey)

	// Find billing that contains the invoice ID
	log.Printf("Searching for billing containing invoice ID: %s", invoiceID)
	billing, err := client.FindBillingByInvoiceID(invoiceID)
	if err != nil {
		return fmt.Errorf("failed to find billing for invoice %s: %w", invoiceID, err)
	}

	log.Printf("Found billing ID: %s (Issue Date: %s, Amount: %d)", billing.ID, billing.IssueDate, billing.Amount)

	// Get download signed URL
	log.Printf("Getting download URL...")
	url, err := client.GetDownloadSignedURL(billing.ID, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	// Download file
	log.Printf("Downloading file...")
	data, err := client.DownloadFile(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Save file
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.pdf", invoiceID))
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("SUCCESS: Saved %s (%d bytes)", filename, len(data))
	log.Println("Download completed successfully")
	return nil
}
