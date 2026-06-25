package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func Rows(rows []map[string]any, format string, orderedColumns ...[]string) error {
	columns := collectColumns(rows, orderedColumns...)
	if format == "json" {
		enc, err := marshalOrderedRowsJSON(rows, columns)
		if err != nil {
			return err
		}
		fmt.Println(string(enc))
		return nil
	}
	if format == "yaml" {
		enc, err := marshalOrderedRowsYAML(rows, columns)
		if err != nil {
			return err
		}
		fmt.Print(string(enc))
		return nil
	}
	if format == "toml" {
		enc, err := marshalOrderedRowsTOML(rows, columns)
		if err != nil {
			return err
		}
		fmt.Print(string(enc))
		return nil
	}
	if len(rows) == 0 {
		fmt.Println("(empty)")
		return nil
	}
	widths := columnWidths(columns, rows)
	t := table.New(
		table.WithColumns(toColumns(columns, widths)),
		table.WithRows(toRows(columns, rows)),
		table.WithHeight(tableHeight(len(rows))),
		table.WithFocused(false),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("8"))
	styles.Selected = styles.Selected.Foreground(lipgloss.Color("15"))
	t.SetStyles(styles)
	fmt.Println(t.View())
	return nil
}

func tableHeight(rowCount int) int {
	const maxHeight = 12
	if rowCount < 1 {
		return 1
	}
	if rowCount > maxHeight {
		return maxHeight
	}
	return rowCount
}

func collectColumns(rows []map[string]any, preferred ...[]string) []string {
	seen := map[string]bool{}
	columns := []string{}
	if len(preferred) > 0 {
		for _, key := range preferred[0] {
			if !seen[key] {
				seen[key] = true
				columns = append(columns, key)
			}
		}
	}
	for _, row := range rows {
		keys := make([]string, 0, len(row))
		for key := range row {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if !seen[key] {
				seen[key] = true
				columns = append(columns, key)
			}
		}
	}
	return columns
}

func marshalOrderedRowsJSON(rows []map[string]any, columns []string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	for rowIndex, row := range rows {
		if rowIndex > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n  {")
		for colIndex, column := range columns {
			if colIndex > 0 {
				buf.WriteString(",")
			}
			key, err := json.Marshal(column)
			if err != nil {
				return nil, err
			}
			value, err := json.Marshal(row[column])
			if err != nil {
				return nil, err
			}
			buf.WriteString("\n    ")
			buf.Write(key)
			buf.WriteString(": ")
			buf.Write(value)
		}
		buf.WriteString("\n  }")
	}
	if len(rows) > 0 {
		buf.WriteString("\n")
	}
	buf.WriteString("]\n")
	return buf.Bytes(), nil
}

func marshalOrderedRowsYAML(rows []map[string]any, columns []string) ([]byte, error) {
	root := yaml.Node{Kind: yaml.SequenceNode}
	for _, row := range rows {
		item := yaml.Node{Kind: yaml.MappingNode}
		for _, column := range columns {
			key := yaml.Node{Kind: yaml.ScalarNode, Value: column}
			value := yaml.Node{}
			if err := value.Encode(row[column]); err != nil {
				return nil, err
			}
			item.Content = append(item.Content, &key, &value)
		}
		root.Content = append(root.Content, &item)
	}
	return yaml.Marshal(&root)
}

func marshalOrderedRowsTOML(rows []map[string]any, columns []string) ([]byte, error) {
	var buf bytes.Buffer
	for rowIndex, row := range rows {
		if rowIndex > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString("[[rows]]\n")
		for _, column := range columns {
			value, ok := row[column]
			if !ok || value == nil {
				continue
			}
			line, err := toml.Marshal(map[string]any{column: value})
			if err != nil {
				return nil, err
			}
			buf.Write(line)
		}
	}
	return buf.Bytes(), nil
}

func columnWidths(columns []string, rows []map[string]any) map[string]int {
	widths := map[string]int{}
	for _, column := range columns {
		widths[column] = len(column)
	}
	for _, row := range rows {
		for _, column := range columns {
			cell := fmt.Sprint(row[column])
			if len(cell) > 48 {
				cell = cell[:48]
			}
			if len(cell) > widths[column] {
				widths[column] = len(cell)
			}
		}
	}
	for column, width := range widths {
		if width < 8 {
			widths[column] = 8
		}
		if width > 64 {
			widths[column] = 64
		}
	}
	return widths
}

func toColumns(names []string, widths map[string]int) []table.Column {
	out := make([]table.Column, 0, len(names))
	for _, name := range names {
		out = append(out, table.Column{Title: strings.ToUpper(name), Width: widths[name]})
	}
	return out
}

func toRows(columns []string, rows []map[string]any) []table.Row {
	out := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		values := make([]string, 0, len(columns))
		for _, column := range columns {
			cell := fmt.Sprint(row[column])
			if len(cell) > 64 {
				cell = cell[:61] + "..."
			}
			values = append(values, cell)
		}
		out = append(out, table.Row(values))
	}
	return out
}
