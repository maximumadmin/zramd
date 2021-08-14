package kernelversion

import "testing"

func TestGte(t *testing.T) {
	tables := []struct {
		a        version
		b        version
		expected bool
	}{
		{version{0, 0}, version{0, 1}, false},
		{version{0, 1}, version{0, 1}, true},
		{version{0, 10}, version{0, 1}, true},
		{version{0, 10}, version{0, 10}, true},
		{version{2, 6}, version{2, 10}, false},
		{version{3, 0}, version{2, 6}, true},
		{version{3, 10}, version{3, 14}, false},
		{version{3, 14}, version{3, 14}, true},
		{version{3, 15}, version{3, 14}, true},
		{version{4, 0}, version{3, 14}, true},
		{version{4, 0}, version{5, 0}, false},
		{version{5, 0}, version{5, 0}, true},
		{version{5, 10}, version{5, 0}, true},
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
