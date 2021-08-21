package utsname

import (
	"strconv"
	"strings"
	"syscall"
)

// UTSName contains information about the current kernel.
type UTSName struct {
	SysName    string
	NodeName   string
	Release    string
	Version    string
	Machine    string
	DomainName string
}

// KernelVersion will split the Release field and return the fist two numbers.
func (uname *UTSName) KernelVersion() (int, int) {
	parts := strings.Split(uname.Release, ".")
	major, _ := strconv.ParseInt(parts[0], 10, strconv.IntSize)
	minor, _ := strconv.ParseInt(parts[1], 10, strconv.IntSize)
	return int(major), int(minor)
}

// Uname returns information about the current kernel.
func Uname() *UTSName {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		panic(err)
	}
	// Keep in mind that we are using 2 sightly different implementations to parse
	// char slices as they use different data types depending on the architecture,
	// see also https://github.com/golang/go/issues/13318.
	return &UTSName{
		SysName:    parseCharSlice(uname.Sysname[:]),
		NodeName:   parseCharSlice(uname.Nodename[:]),
		Release:    parseCharSlice(uname.Release[:]),
		Version:    parseCharSlice(uname.Version[:]),
		Machine:    parseCharSlice(uname.Machine[:]),
		DomainName: parseCharSlice(uname.Domainname[:]),
	}
}
