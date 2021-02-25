package memory

import (
	"testing"
)

func TestParseMemInfoLine(t *testing.T) {
	tables := []struct {
		input string
		key   string
		value uint64
	}{
		{"MemTotal:       65790200 kB", "MemTotal", 65790200},
		{"Inactive(file):  2770308 kB", "Inactive(file)", 2770308},
		{"VmallocTotal:   54358348362 kB", "VmallocTotal", 54358348362},
		{"ShmemHugePages:        0 kB", "ShmemHugePages", 0},
		{"HugePages_Surp:        0", "HugePages_Surp", 0},
		{"Hugepagesize:       2048 kB", "Hugepagesize", 2048},
	}
	for _, table := range tables {
		key, value := parseMemInfoLine(table.input)
		if key != table.key {
			t.Errorf("expected \"%s\" but got \"%s\" instead", table.key, key)
		}
		if value != table.value {
			t.Errorf("expected \"%d\" but got \"%d\" instead", table.value, value)
		}
	}
}
