package zram

import (
	"fmt"
	"zramd/src/util"
)

// MakeSwap formats a zram device given an index corresponding to the zram
// device path under /dev.
func MakeSwap(index int) error {
	file := fmt.Sprintf("/dev/zram%d", index)
	return util.Run("mkswap", file)
}

// SwapOn enables a swap device given a zram device index and a priority.
func SwapOn(index int, priority int) error {
	file := fmt.Sprintf("/dev/zram%d", index)
	return util.Run("swapon", file, "--priority", fmt.Sprint(priority))
}

// SwapOff disables a swap device given a zram device index.
func SwapOff(index int) error {
	file := fmt.Sprintf("/dev/zram%d", index)
	return util.Run("swapoff", file)
}
