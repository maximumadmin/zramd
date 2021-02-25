package kernelversion

import "zramd/src/utsname"

var major, minor = utsname.Uname().KernelVersion()

// SupportsZram checks if current kernel version supports zram.
func SupportsZram() bool {
	return (major == 3 && minor >= 14) || major > 3
}

// SupportsZstd checks if current kernel supports zstd compressed zram.
func SupportsZstd() bool {
	return (major == 4 && minor >= 19) || major > 4
}
