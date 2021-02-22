package zram

import (
	"fmt"
	"os"
)

// DeviceExists checks if a zram device exists.
func DeviceExists(index int) bool {
	_, err := os.Stat(fmt.Sprintf("/dev/zram%d", index))
	if err != nil {
		return false
	}
	return true
}

func setAttribute(index int, name string, value string) error {
	file := fmt.Sprintf("/sys/block/zram%d/%s", index, name)
	data := []byte(value)
	return os.WriteFile(file, data, 0644)
}

// SetSize sets the size in bytes for a zram device.
func SetSize(index int, size uint64) error {
	return setAttribute(index, "disksize", fmt.Sprint(size))
}

// SetCompAlgorithm sets the compression algorithm for a zram device.
func SetCompAlgorithm(index int, algorithm string) error {
	return setAttribute(index, "comp_algorithm", algorithm)
}
