package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func fmtTableOutput(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	//table.EnableBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

func loadHistory(limit int) (History, error) {
	var history History

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("Error homedir history file: %s\n", err)
	}
	file, err := os.Open(filepath.Join(home, ".ots_history"))

	if err != nil {
		return nil, fmt.Errorf("opening history file: %v", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into timestamp and key
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %v", err)
		}

		entry := strings.TrimSpace(parts[1])
		if len(entry) == 0 {
			continue
		}
		history = append(history, entry)
		// to do check expired entries to not load
		// timeNow := time.Now()
		// strTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
		// if err != nil {
		// 	fmt.Printf("Error converting string to int: %v\n", err)
		// }
		// timestamp := time.Unix(strTimestamp, 0)

		// expired := timeNow.Before(timestamp)

		// if !expired {
		// 	history = append(history, entry)
		// }

	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading history file: %v", err)
	}

	if len(history) >= limit {
		history = history[len(history)-limit:]
	}

	return history, nil

}

func writeHistory(entry string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting homedir: %s\n", err)
	}

	file, err := os.OpenFile(filepath.Join(home, ".ots_history"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening history file: %s", err)
	}
	defer file.Close()

	// Format the entry as a string
	entryString := fmt.Sprintf("%d;%s\n", time.Now().Unix(), entry)

	if _, err := file.WriteString(entryString); err != nil {
		return fmt.Errorf("error writing entry to file: %s", err)
	}

	return nil
}
