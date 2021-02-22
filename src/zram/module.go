package zram

import (
	"fmt"
	"os"
	"strings"
	"zramd/src/util"
)

// LoadModule loads the zram module.
func LoadModule(n int) error {
	return util.Run("modprobe", "zram", fmt.Sprintf("num_devices=%d", n))
}

// UnloadModule unloads the zram module.
func UnloadModule() error {
	return util.Run("modprobe", "-r", "zram")
}

// IsLoaded checks if the zram module has been loaded.
func IsLoaded() bool {
	// Reading from /proc/modules should be faster than using lsmod.
	data, err := os.ReadFile("/proc/modules")
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "zram ") {
			return true
		}
	}
	return false
}
