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

func OutputBulkTable(h []string, d [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(h)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(d) // Add Bulk Data
	table.Render()
}

func OutputTable(h []string, d []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(h)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.Append(d)
	table.Render()
}

func LoadHistory(limit int) (History, error) {
	var history History

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("Error homedir history file: %s\n", err)
	}
	file, err := os.Open(filepath.Join(home, ".ots_history"))

	if err != nil {
		fmt.Printf("Error opening history file: %s\n", err)
		return nil, err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into timestamp and key
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			fmt.Printf("Invalid line format: %s\n", line)
			continue
		}
		entry := strings.TrimSpace(parts[1])
		history = append(history, entry)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading history file: %s\n", err)
		return nil, err
	}
	if len(history) >= limit {
		history = history[len(history)-limit:]
	}

	return history, nil

}

func WriteHistory(entry string) error {
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
