package kernelversion

import "testing"

func TestGte(t *testing.T) {
	tables := []struct {
		a        kVersion
		b        kVersion
		expected bool
	}{
		{kVersion{0, 0}, kVersion{0, 1}, false},
		{kVersion{0, 1}, kVersion{0, 1}, true},
		{kVersion{0, 10}, kVersion{0, 1}, true},
		{kVersion{0, 10}, kVersion{0, 10}, true},
		{kVersion{2, 6}, kVersion{2, 10}, false},
		{kVersion{3, 0}, kVersion{2, 6}, true},
		{kVersion{3, 10}, kVersion{3, 14}, false},
		{kVersion{3, 14}, kVersion{3, 14}, true},
		{kVersion{3, 15}, kVersion{3, 14}, true},
		{kVersion{4, 0}, kVersion{3, 14}, true},
		{kVersion{4, 0}, kVersion{5, 0}, false},
		{kVersion{5, 0}, kVersion{5, 0}, true},
		{kVersion{5, 10}, kVersion{5, 0}, true},
	}
	for _, table := range tables {
		if table.expected != gte(table.a, table.b) {
			t.Errorf(
				"wrong comparison result for %d.%d >= %d.%d",
				table.a.major,
				table.a.minor,
				table.b.major,
				table.b.minor,
			)
		}
	}
}
