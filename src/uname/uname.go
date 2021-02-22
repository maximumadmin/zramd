package uname

import (
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
