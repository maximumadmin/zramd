package zram

import (
	"fmt"
)

// MakeSwap formats a zram device given an index corresponding to the zram
// device path under "/dev", this process is very fast and there is no
// noticeable delay if ran multiple times sequentially.
func MakeSwap(index int) error {
	file := fmt.Sprintf("/dev/zram%d", index)
	return execute("mkswap", file)
}

// SwapOn enables a swap device given a zram device index and a priority, this
// process is slow (about 60ms per swap device on a 16-core CPU, depends on the
// swap size and hardware), specially with large and multiple swap devices.
func SwapOn(index int, priority int) error {
	file := fmt.Sprintf("/dev/zram%d", index)
	return execute("swapon", file, "--priority", fmt.Sprint(priority))
}

// SwapOff disables a swap device given a zram device index.
func SwapOff(index int) error {
	file := fmt.Sprintf("/dev/zram%d", index)
	return execute("swapoff", file)
}
