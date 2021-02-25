package zram

import (
	"fmt"
	"os"
	"strings"
)

// LoadModule loads the zram module, we need to specify the number of devices
// when loading the zram module, otherwise only one zram device will be created,
// you can verify the number of zram devices by running "ls /dev", zram devices
// start at zram0.
func LoadModule(n int) error {
	return execute("modprobe", "zram", fmt.Sprintf("num_devices=%d", n))
}

// UnloadModule unloads the zram module.
func UnloadModule() error {
	return execute("modprobe", "-r", "zram")
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
