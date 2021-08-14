package zram

import (
	"fmt"
	"os"
)

// DeviceExists checks if a zram device exists.
func DeviceExists(id int) bool {
	_, err := os.Stat(fmt.Sprintf("/dev/zram%d", id))
	return err == nil
}

func setAttribute(id int, name string, value string) error {
	file := fmt.Sprintf("/sys/block/zram%d/%s", id, name)
	data := []byte(value)
	return os.WriteFile(file, data, 0644)
}

// setSize sets the size in bytes for a zram device.
func setSize(id int, size uint64) error {
	return setAttribute(id, "disksize", fmt.Sprint(size))
}

// setCompAlgorithm sets the compression algorithm for a zram device.
func setCompAlgorithm(id int, algorithm string) error {
	return setAttribute(id, "comp_algorithm", algorithm)
}

// Configure sets the size and compression algorithm of a zram device, see also
// https://www.kernel.org/doc/html/latest/admin-guide/blockdev/zram.html#deactivate.
func Configure(id int, size uint64, algorithm string) error {
	// Compression algorithm must be set before size, otherwise you will get
	// errors like "device or resource busy".
	if err := setCompAlgorithm(id, algorithm); err != nil {
		return err
	}
	return setSize(id, size)
}
