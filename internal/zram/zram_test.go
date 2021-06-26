package zram

import "testing"

func TestGetZramID(t *testing.T) {
	tables := []struct {
		input    string
		expected int
	}{
		{"Filename				Type		Size		Used		Priority", -1},
		{"/zram324                                 partition	262140		0		100", 324},
		{"/zram16                                 partition	262140		0		100", 16},
		{"/zram0                                 partition	262140		0		100", 0},
		{"/dev/zram324                                 partition	262140		0		100", 324},
		{"/dev/zram16                                 partition	262140		0		100", 16},
		{"/dev/zram0                                 partition	262140		0		100", 0},
		{"/zram18                               file		524284		0		-2", -1},
		{"/zram                               file		524284		0		-2", -1},
		{"/swapfile                               file		524284		0		-2", -1},
	}
	for _, table := range tables {
		got := getZramID(table.input)
		if table.expected != got {
			t.Errorf("expected \"%d\" but got \"%d\" instead", table.expected, got)
		}
	}
}
