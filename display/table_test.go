package display

import (
	"strings"
	"testing"
)

func TestTableHeight(t *testing.T) {
	tests := []struct {
		name string
		rows int
		want int
	}{
		{name: "empty", rows: 0, want: 1},
		{name: "one row", rows: 1, want: 1},
		{name: "few rows", rows: 3, want: 3},
		{name: "capped", rows: 20, want: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tableHeight(tt.rows)
			if got != tt.want {
				t.Fatalf("tableHeight(%d) = %d, want %d", tt.rows, got, tt.want)
			}
		})
	}
}

func TestCollectColumnsUsesPreferredOrder(t *testing.T) {
	rows := []map[string]any{
		{"b": 2, "a": 1, "c": 3},
	}
	got := collectColumns(rows, []string{"c", "a"})
	want := []string{"c", "a", "b"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("columns = %v, want %v", got, want)
	}
}

func TestMarshalOrderedRowsJSON(t *testing.T) {
	rows := []map[string]any{{"b": 2, "a": 1}}
	got, err := marshalOrderedRowsJSON(rows, []string{"b", "a"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "\"b\": 2,\n    \"a\": 1") {
		t.Fatalf("JSON did not preserve column order:\n%s", got)
	}
}

func TestMarshalOrderedRowsYAML(t *testing.T) {
	rows := []map[string]any{{"b": 2, "a": 1}}
	got, err := marshalOrderedRowsYAML(rows, []string{"b", "a"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "b: 2\n  a: 1") {
		t.Fatalf("YAML did not preserve column order:\n%s", got)
	}
}

func TestMarshalOrderedRowsTOML(t *testing.T) {
	rows := []map[string]any{{"b": 2, "a": 1}}
	got, err := marshalOrderedRowsTOML(rows, []string{"b", "a"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "[[rows]]\nb = 2\na = 1") {
		t.Fatalf("TOML did not preserve column order:\n%s", got)
	}
}
