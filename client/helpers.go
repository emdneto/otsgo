package client

import (
	"os"

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
	//
}
