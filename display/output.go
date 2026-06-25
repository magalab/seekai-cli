package display

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func ResolveFormat(value, fallback string) string {
	if value == "" || value == "auto" {
		return fallback
	}
	return value
}

func IsStructured(format string) bool {
	switch format {
	case "json", "yaml", "toml":
		return true
	default:
		return false
	}
}

func Value(value any, format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(value)
	case "yaml":
		data, err := yaml.Marshal(value)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	case "toml":
		data, err := toml.Marshal(tomlDocument(value))
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	case "table":
		return Rows(valueRows(value), "table")
	case "text", "auto":
		switch v := value.(type) {
		case string:
			fmt.Println(v)
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

func valueRows(value any) []map[string]any {
	switch v := value.(type) {
	case map[string]any:
		return []map[string]any{v}
	case map[string]string:
		row := make(map[string]any, len(v))
		for key, value := range v {
			row[key] = value
		}
		return []map[string]any{row}
	default:
		return []map[string]any{{"value": value}}
	}
}

func tomlDocument(value any) any {
	switch value.(type) {
	case map[string]any, map[string]string, map[string]int, map[string]float64, map[string]bool:
		return value
	default:
		return map[string]any{"value": value}
	}
}

func LongText(text string, w io.Writer) error {
	if text == "" {
		_, err := fmt.Fprintln(w, "(empty)")
		return err
	}
	if file, ok := w.(*os.File); ok && file == os.Stdout && isTerminal(file) {
		return RunViewport(text)
	}
	_, err := fmt.Fprintln(w, text)
	return err
}

func isTerminal(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice != 0
}
