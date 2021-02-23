package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"zramd/src/memory"
	"zramd/src/util"
	"zramd/src/zram"

	"github.com/alexflint/go-arg"
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

func swapOn(index int, priority int, c chan error) {
	if err := zram.MakeSwap(index); err != nil {
		c <- err
		return
	}
	if err := zram.SwapOn(index, priority); err != nil {
		c <- err
		return
	}
	c <- nil
}

func setupSwap(numCPU int, swapPriority int) []error {
	var wg sync.WaitGroup
	var errors []error
	channel := make(chan error)
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			swapOn(index, swapPriority, channel)
		}(i)
	}
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
	if zram.IsLoaded() {
		errorf("the zram module is already loaded")
		return 1
	}
	numCPU := runtime.NumCPU()
	if err := zram.LoadModule(numCPU); err != nil {
		errorf(err.Error())
		return 1
	}
	maxTotalSize := getMaxTotalSize(
		memory.ReadMemInfo()["MemTotal"]*1024,
		uint64(cmd.MaxSizeMB)*1024*1024,
		cmd.MaxSizePercent,
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
	if !zram.IsLoaded() {
		errorf("the zram module is not loaded")
		return 1
	}
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

func run() int {
	var args args
	parser := arg.MustParse(&args)
	if parser.Subcommand() == nil {
		parser.Fail("missing subcommand")
	}

	switch {
	case args.Start != nil:
		if args.Start.Algorithm == "zstd" && !util.IsZstdSupported() {
			parser.Fail("the zstd algorithm is not supported on kernels < 4.19")
		}
		if args.Start.MaxSizePercent < 0.05 || args.Start.MaxSizePercent > 1 {
			parser.Fail("--max-ram must be a value between 0.05 and 1")
		}
		if !util.IsRoot() {
			errorf("root privileges are required")
			return 1
		}
		return initializeZram(args.Start)

	case args.Stop != nil:
		if !util.IsRoot() {
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
