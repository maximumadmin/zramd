package zram

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MakeSwap formats a zram device given a zram device id, this process is very
// fast and there is no noticeable delay if ran multiple times sequentially.
func MakeSwap(id int) error {
	file := fmt.Sprintf("/dev/zram%d", id)
	return execute("mkswap", file)
}

// SwapOn enables a swap device given a zram device id and a priority, this
// process is slow (about 60ms per swap device on a 16-core CPU, depends on the
// swap size and hardware), specially with large and multiple swap devices.
func SwapOn(id int, priority int) error {
	file := fmt.Sprintf("/dev/zram%d", id)
	return execute("swapon", file, "--priority", fmt.Sprint(priority))
}

// SwapOff disables a swap device given a zram device id.
func SwapOff(id int) error {
	file := fmt.Sprintf("/dev/zram%d", id)
	return execute("swapoff", file)
}

// getZramID parses lines like "/zram16 partition 262140 0 100" or
// "/swapfile file 524284 0 -2" and returns the zram device id if filename (1st
// column) belongs to a zram device and type (2nd column) is "partition" (so we
// avoid cases when the user has a swap file called "zram" 🤦‍♂️), if line does not
// match the previous conditions, the returned value will be -1.
func getZramID(line string) int {
	fields := strings.Fields(line)
	// We need at least 2 columns (filename and type) and type to be partition.
	if len(fields) < 2 || fields[1] != "partition" {
		return -1
	}
	// Depending of how the zram device was created (inside or outside of the
	// systemd sandbox), the filename can sometimes be "/dev/zramX" or "/zramX",
	// so split the filename to get the last part only.
	parts := strings.Split(fields[0], "/")
	filename := parts[len(parts)-1]
	if !strings.HasPrefix(filename, "zram") {
		return -1
	}
	id, err := strconv.ParseInt(
		strings.TrimPrefix(filename, "zram"),
		10,
		strconv.IntSize,
	)
	if err != nil {
		return -1
	}
	return int(id)
}

// SwapDeviceIDs returns a list of the zram device IDs currently used as swap.
func SwapDeviceIDs() []int {
	data, err := os.ReadFile("/proc/swaps")
	if err != nil {
		panic(err)
	}
	result := []int{}
	for _, line := range strings.Split(string(data), "\n") {
		id := getZramID(line)
		if id > -1 {
			result = append(result, id)
		}
	}
	return result
}
