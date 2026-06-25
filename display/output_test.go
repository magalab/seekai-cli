package display

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestIsStructured(t *testing.T) {
	for _, format := range []string{"json", "yaml", "toml"} {
		if !IsStructured(format) {
			t.Fatalf("IsStructured(%q) = false, want true", format)
		}
	}
	for _, format := range []string{"auto", "table", "text"} {
		if IsStructured(format) {
			t.Fatalf("IsStructured(%q) = true, want false", format)
		}
	}
}

func TestTOMLDocumentWrapsScalars(t *testing.T) {
	doc, ok := tomlDocument("hello").(map[string]any)
	if !ok {
		t.Fatalf("tomlDocument returned %T, want map[string]any", doc)
	}
	if doc["value"] != "hello" {
		t.Fatalf("value = %v, want hello", doc["value"])
	}
}

func TestTOMLRowsDocumentEncodes(t *testing.T) {
	rows := map[string]any{
		"rows": []map[string]any{
			{"name": "demo", "type": "completion"},
		},
	}
	if _, err := toml.Marshal(rows); err != nil {
		t.Fatalf("toml.Marshal(rows) error = %v", err)
	}
}

func TestValueRows(t *testing.T) {
	rows := valueRows(map[string]any{"status": "ok"})
	if len(rows) != 1 || rows[0]["status"] != "ok" {
		t.Fatalf("valueRows(map) = %#v", rows)
	}

	rows = valueRows("ok")
	if len(rows) != 1 || rows[0]["value"] != "ok" {
		t.Fatalf("valueRows(scalar) = %#v", rows)
	}
}
