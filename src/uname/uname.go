package uname

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

func parseInt8(data []int8) string {
	b := make([]byte, 0, len(data))
	for _, v := range data {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

// Uname returns information about the current kernel.
func Uname() *UTSName {
	var uname syscall.Utsname
	err := syscall.Uname(&uname)
	if err != nil {
		panic(err)
	}
	return &UTSName{
		SysName:    parseInt8(uname.Sysname[:]),
		NodeName:   parseInt8(uname.Nodename[:]),
		Release:    parseInt8(uname.Release[:]),
		Version:    parseInt8(uname.Version[:]),
		Machine:    parseInt8(uname.Machine[:]),
		DomainName: parseInt8(uname.Domainname[:]),
	}
}
