package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func Print(data json.RawMessage, format string) error {
	if data == nil {
		return nil
	}

	switch format {
	case "json", "":
		return printJSON(data, os.Stdout)
	case "table":
		return printTable(data, os.Stdout)
	case "csv":
		return printCSV(data, os.Stdout)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func printJSON(data json.RawMessage, w io.Writer) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func printTable(data json.RawMessage, w io.Writer) error {
	rows := extractRows(data)
	if len(rows) == 0 {
		fmt.Fprintln(w, "(no data)")
		return nil
	}

	keys := extractKeys(rows)
	if len(keys) == 0 {
		return printJSON(data, w)
	}

	widths := make([]int, len(keys))
	for i, k := range keys {
		widths[i] = len(k)
	}
	stringRows := make([][]string, len(rows))
	for i, row := range rows {
		stringRows[i] = make([]string, len(keys))
		for j, k := range keys {
			val := formatValue(row[k])
			stringRows[i][j] = val
			if len(val) > widths[j] {
				widths[j] = len(val)
			}
		}
	}

	for i, w := range widths {
		if w > 50 {
			widths[i] = 50
		}
	}

	printRow(w, keys, widths)
	seps := make([]string, len(keys))
	for i, width := range widths {
		seps[i] = strings.Repeat("-", width)
	}
	printRow(w, seps, widths)

	for _, row := range stringRows {
		printRow(w, row, widths)
	}

	return nil
}

func printRow(w io.Writer, cols []string, widths []int) {
	parts := make([]string, len(cols))
	for i, col := range cols {
		if len(col) > widths[i] {
			col = col[:widths[i]-2] + ".."
		}
		parts[i] = fmt.Sprintf("%-*s", widths[i], col)
	}
	fmt.Fprintln(w, strings.Join(parts, "  "))
}

func printCSV(data json.RawMessage, w io.Writer) error {
	rows := extractRows(data)
	if len(rows) == 0 {
		return nil
	}

	keys := extractKeys(rows)
	if len(keys) == 0 {
		return printJSON(data, w)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write(keys); err != nil {
		return err
	}
	for _, row := range rows {
		record := make([]string, len(keys))
		for i, k := range keys {
			record[i] = formatValue(row[k])
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func extractRows(data json.RawMessage) []map[string]interface{} {
	var list struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(data, &list); err == nil && list.Items != nil {
		return list.Items
	}

	var arr []map[string]interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr
	}

	var single map[string]interface{}
	if err := json.Unmarshal(data, &single); err == nil {
		return []map[string]interface{}{single}
	}

	return nil
}

func extractKeys(rows []map[string]interface{}) []string {
	seen := make(map[string]bool)
	var keys []string
	for _, row := range rows {
		for k := range row {
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	return keys
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
