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

// SupportsMultiCompStreams checks if current kernel supports multiple
// compression streams, this feature is required in order to take advantage of
// multiple processors with a single zram device, see also
// https://wiki.gentoo.org/wiki/Zram#Caveats.2FCons.
func SupportsMultiCompStreams() bool {
	return (major == 3 && minor >= 15) || major > 3
}
