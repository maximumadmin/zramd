package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"zramd/src/kernelversion"
	"zramd/src/memory"
	"zramd/src/zram"

	"github.com/alexflint/go-arg"
)

// startCmd contains the arguments used by the start subcommand, Fraction will
// be the same size as the physical memory by default, see also
// https://fedoraproject.org/wiki/Changes/Scale_ZRAM_to_full_memory_size.
type startCmd struct {
	Algorithm    string  `arg:"-a,env" default:"zstd" help:"zram compression algorithm"`
	Fraction     float32 `arg:"-f,env" default:"1.0" help:"maximum percentage of RAM allowed to use"`
	MaxSizeMB    int     `arg:"-m,--max-size,env:MAX_SIZE" default:"8192" placeholder:"MAX_SIZE" help:"maximum total MB of swap to allocate"`
	SwapPriority int     `arg:"-p,--priority,env:PRIORITY" default:"100" placeholder:"PRIORITY" help:"swap priority"`
}

type stopCmd struct {
}

type args struct {
	Start *startCmd `arg:"subcommand:start" help:"load zram module and setup swap devices"`
	Stop  *stopCmd  `arg:"subcommand:stop" help:"stop swap devices and unload zram module"`
}

func errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
}

// getMaxTotalSize gets the maximum amount of memory (in bytes) that is going to
// be used for the creation of the swap-on-RAM devices.
func getMaxTotalSize(
	memTotalBytes uint64,
	maxSizeBytes uint64,
	maxPercent float32,
) uint64 {
	size := uint64(float32(memTotalBytes) * maxPercent)
	if size < maxSizeBytes {
		return size
	}
	return maxSizeBytes
}

func swapOn(id int, priority int, c chan error) {
	if err := zram.MakeSwap(id); err != nil {
		c <- err
		return
	}
	if err := zram.SwapOn(id, priority); err != nil {
		c <- err
		return
	}
	c <- nil
}

// setupSwap will initialize the swap devices in parallel, this operation will
// not make swap initialization numCPU times faster, but it will still be faster
// than doing it sequentially.
func setupSwap(numCPU int, swapPriority int) []error {
	var wg sync.WaitGroup
	var errors []error
	channel := make(chan error)
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			swapOn(id, swapPriority, channel)
		}(i)
	}
	// Using a separate routine to extract data from the channel, see also
	// https://stackoverflow.com/a/54535532.
	go func() {
		for err := range channel {
			if err != nil {
				errors = append(errors, err)
			}
		}
	}()
	wg.Wait()
	close(channel)
	return errors
}

func initializeZram(cmd *startCmd) int {
	numCPU := runtime.NumCPU()
	if err := zram.LoadModule(numCPU); err != nil {
		errorf(err.Error())
		return 1
	}
	maxTotalSize := getMaxTotalSize(
		memory.ReadMemInfo()["MemTotal"]*1024,
		uint64(cmd.MaxSizeMB)*1024*1024,
		cmd.Fraction,
	)
	deviceSize := maxTotalSize / uint64(numCPU)
	for i := 0; i < numCPU; i++ {
		if !zram.DeviceExists(i) {
			errorf("device zram%d does not exist", i)
			return 1
		}
		err := zram.Configure(i, deviceSize, cmd.Algorithm)
		if err != nil {
			errorf(err.Error())
			return 1
		}
	}
	errors := setupSwap(numCPU, cmd.SwapPriority)
	if len(errors) > 0 {
		for _, err := range errors {
			errorf(err.Error())
		}
		return 1
	}
	return 0
}

func deinitializeZram() int {
	ret := 0
	for i := 0; i < runtime.NumCPU(); i++ {
		if !zram.DeviceExists(i) {
			continue
		}
		if err := zram.SwapOff(i); err != nil {
			errorf("zram%d: %s", i, err.Error())
			ret = 1
		}
	}
	if err := zram.UnloadModule(); err != nil {
		errorf(err.Error())
		ret = 1
	}
	return ret
}

func isRoot() bool {
	return os.Geteuid() == 0
}

func run() int {
	if !kernelversion.SupportsZram() {
		errorf("zram is not supported on kernels < 3.14")
		return 1
	}

	var args args
	parser := arg.MustParse(&args)
	if parser.Subcommand() == nil {
		parser.Fail("missing subcommand")
	}

	switch {
	case args.Start != nil:
		if args.Start.Algorithm == "zstd" && !kernelversion.SupportsZstd() {
			parser.Fail("the zstd algorithm is not supported on kernels < 4.19")
		}
		if args.Start.Fraction < 0.05 || args.Start.Fraction > 1 {
			parser.Fail("--fraction must have a value between 0.05 and 1")
		}
		if args.Start.SwapPriority < -1 || args.Start.SwapPriority > 32767 {
			parser.Fail("--priority must have a value between -1 and 32767")
		}
		if zram.IsLoaded() {
			errorf("the zram module is already loaded")
			return 1
		}
		if !isRoot() {
			errorf("root privileges are required")
			return 1
		}
		return initializeZram(args.Start)

	case args.Stop != nil:
		if !zram.IsLoaded() {
			errorf("the zram module is not loaded")
			return 1
		}
		if !isRoot() {
			errorf("root privileges are required")
			return 1
		}
		return deinitializeZram()
	}

	return 0
}

func main() {
	code := run()
	os.Exit(code)
}
