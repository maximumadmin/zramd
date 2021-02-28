package kernelversion

import "zramd/src/utsname"

var major, minor = utsname.Uname().KernelVersion()

type kVersion struct {
	major int
	minor int
}

func gte(a kVersion, b kVersion) bool {
	return (a.major == b.major && a.minor >= b.minor) || a.major > b.major
}

// SupportsZram checks if current kernel version supports zram.
func SupportsZram() bool {
	return gte(kVersion{major, minor}, kVersion{3, 14})
}

// SupportsZstd checks if current kernel supports zstd compressed zram.
func SupportsZstd() bool {
	return gte(kVersion{major, minor}, kVersion{4, 19})
}

// SupportsMultiCompStreams checks if current kernel supports multiple
// compression streams, this feature is required in order to take advantage of
// multiple processors with a single zram device, see also
// https://wiki.gentoo.org/wiki/Zram#Caveats.2FCons.
func SupportsMultiCompStreams() bool {
	return gte(kVersion{major, minor}, kVersion{3, 15})
}
