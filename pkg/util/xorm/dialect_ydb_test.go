package xorm

import "testing"

func TestGetLastPartOfColumn(t *testing.T) {
	tests := []struct {
		column   string
		expected string
	}{
		{"table.column", "column"},
		{"column", "column"},
		{".", ""},
		{".x", "x"},
		{"xxxx.xxx.", ""},
	}

	for _, test := range tests {
		if getLastPartOfColumn(test.column) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, getLastPartOfColumn(test.column))
		}
	}
}
