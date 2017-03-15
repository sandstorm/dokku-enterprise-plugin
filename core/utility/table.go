package utility

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
	"fmt"
)

func RenderAsTable(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}

func ParseTable(table string) (rows [][]string, err error){
	tableRows := strings.Split(table, "\n")

	rowsWithContent := filterContentRows(tableRows)

	if len(rowsWithContent) == 0 {
		return nil, fmt.Errorf("Table does not contain any valid content. Tried to parse: %v", table)
	}

	for _, row := range rowsWithContent {
		rows = append(rows, parseTableRow(row))
	}

	return
}

func filterContentRows(rows []string) (result []string) {
	for _, row := range rows {
		if strings.HasPrefix(row, "|") {
			result = append(result, row)
		}
	}

	return
}

func parseTableRow(row string) (result []string) {
	cells := strings.Split(row, "|")

	for _, cell := range cells[1:len(cells)-1] {
		result = append(result, strings.TrimSpace(cell))
	}

	return
}