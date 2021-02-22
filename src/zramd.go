package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"

	"zramd/src/uname"
)

type startCmd struct {
	Algorithm      string  `arg:"-a,env" default:"zstd" placeholder:"A" help:"zram compression algorithm"`
	MaxSizeMB      int     `arg:"-m,--max-size,env" default:"8192" placeholder:"M" help:"maximum total MB of swap to allocate"`
	MaxSizePercent float32 `arg:"-r,--max-ram,env" default:"0.5" placeholder:"P" help:"maximum percentage of RAM allowed to use"`
	SwapPriority   int     `arg:"-p,--priority,env" default:"10" placeholder:"N" help:"swap priority"`
}

type stopCmd struct {
}

type args struct {
	Start *startCmd `arg:"subcommand:start" help:"load zram module and setup swap devices"`
	Stop  *stopCmd  `arg:"subcommand:stop" help:"stop swap devices and unload zram module"`
}

func parseKernelVersion(version string) (int, int) {
	parts := strings.Split(version, ".")
	major, _ := strconv.ParseInt(parts[0], 10, strconv.IntSize)
	minor, _ := strconv.ParseInt(parts[1], 10, strconv.IntSize)
	return int(major), int(minor)
}

func isZstdSupported() bool {
	major, minor := parseKernelVersion(uname.Uname().Release)
	return (major == 4 && minor >= 19) || major > 4
}

func isRoot() bool {
	return os.Geteuid() == 0
}

func run() int {
	var args args
	parser := arg.MustParse(&args)

	switch {
	case args.Start != nil:
		if args.Start.Algorithm == "zstd" && !isZstdSupported() {
			parser.Fail("The zstd algorithm is not supported on kernels < 4.19")
		}
		if !isRoot() {
			fmt.Fprintf(os.Stderr, "Root privileges are required\n")
			return 1
		}
		return 0

	case args.Stop != nil:
		if !isRoot() {
			fmt.Fprintf(os.Stderr, "Root privileges are required\n")
			return 1
		}
		return 0
	}

	return 0
}

func main() {
	code := run()
	os.Exit(code)
}
