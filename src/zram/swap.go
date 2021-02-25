package zram

import (
	"fmt"
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
