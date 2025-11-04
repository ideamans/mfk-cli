package downloader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/miyanaga/moneyforward-kessai-invoice-downloader-v2/internal/mfkessai"
)

// Downloader handles concurrent invoice downloads
type Downloader struct {
	client      *mfkessai.Client
	concurrency int
}

// New creates a new Downloader
func New(client *mfkessai.Client, concurrency int) *Downloader {
	return &Downloader{
		client:      client,
		concurrency: concurrency,
	}
}

// DownloadInvoices downloads all invoices within the specified date range
func (d *Downloader) DownloadInvoices(startDate, endDate, outputDir string) error {
	// Fetch billings
	log.Printf("Fetching billings from %s to %s...", startDate, endDate)
	billings, err := d.client.GetBillings(startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to get billings: %w", err)
	}

	if len(billings) == 0 {
		log.Println("No billings found for the specified date range")
		return nil
	}

	// Count total invoices
	totalInvoices := 0
	for _, b := range billings {
		totalInvoices += len(b.InvoiceIDs)
	}

	log.Printf("Found %d billing(s) with %d invoice(s)", len(billings), totalInvoices)

	if totalInvoices == 0 {
		log.Println("No invoices to download")
		return nil
	}

	// Create a semaphore to limit concurrency
	sem := make(chan struct{}, d.concurrency)
	var wg sync.WaitGroup
	var errCount atomic.Int32
	var successCount atomic.Int32

	// Create error channel with buffer for all possible invoices
	errChan := make(chan error, totalInvoices)

	invoiceNum := 0
	for _, billing := range billings {
		for _, invoiceID := range billing.InvoiceIDs {
			invoiceNum++
			wg.Add(1)
			go func(num int, b mfkessai.Billing, invID string) {
				defer wg.Done()

				// Acquire semaphore
				sem <- struct{}{}
				defer func() { <-sem }()

				log.Printf("[%d/%d] Processing invoice ID: %s (Billing: %s, Issue Date: %s, Amount: %d)",
					num, totalInvoices, invID, b.ID, b.IssueDate, b.Amount)

				// Get download signed URL
				url, err := d.client.GetDownloadSignedURL(b.ID, invID)
				if err != nil {
					errCount.Add(1)
					errChan <- fmt.Errorf("failed to get download URL for invoice %s: %w", invID, err)
					log.Printf("[%d/%d] ERROR: Failed to get download URL for %s: %v",
						num, totalInvoices, invID, err)
					return
				}

				// Download file
				data, err := d.client.DownloadFile(url)
				if err != nil {
					errCount.Add(1)
					errChan <- fmt.Errorf("failed to download file for invoice %s: %w", invID, err)
					log.Printf("[%d/%d] ERROR: Failed to download file for %s: %v",
						num, totalInvoices, invID, err)
					return
				}

				// Save file
				filename := filepath.Join(outputDir, fmt.Sprintf("%s.pdf", invID))
				if err := os.WriteFile(filename, data, 0644); err != nil {
					errCount.Add(1)
					errChan <- fmt.Errorf("failed to write file for invoice %s: %w", invID, err)
					log.Printf("[%d/%d] ERROR: Failed to write file for %s: %v",
						num, totalInvoices, invID, err)
					return
				}

				successCount.Add(1)
				log.Printf("[%d/%d] SUCCESS: Saved %s (%d bytes)",
					num, totalInvoices, filename, len(data))
			}(invoiceNum, billing, invoiceID)
		}
	}

	// Wait for all downloads to complete
	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// Print summary
	log.Printf("\n=== Download Summary ===")
	log.Printf("Total billings: %d", len(billings))
	log.Printf("Total invoices: %d", totalInvoices)
	log.Printf("Successful downloads: %d", successCount.Load())
	log.Printf("Failed downloads: %d", errCount.Load())

	if len(errors) > 0 {
		log.Printf("\nErrors encountered:")
		for i, err := range errors {
			log.Printf("  %d. %v", i+1, err)
		}
		return fmt.Errorf("%d out of %d invoice downloads failed", errCount.Load(), totalInvoices)
	}

	return nil
}
