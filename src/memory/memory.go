package memory

import (
	"os"
	"strconv"
	"strings"
)

// parseMemInfoLine parses lines like "MemFree: 463 kB" or "HugePages_Total: 0",
// keep in mind that a single line can contain multiple whitespaces.
func parseMemInfoLine(line string) (string, uint64) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return "", 0
	}
	value, _ := strconv.ParseUint(fields[1], 10, 64)
	key := strings.TrimSuffix(fields[0], ":")
	return key, uint64(value)
}

// ReadMemInfo reads the values of /proc/meminfo (they will always be in KiB),
// see also https://unix.stackexchange.com/a/199491.
func ReadMemInfo() map[string]uint64 {
	result := map[string]uint64{}
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		key, value := parseMemInfoLine(line)
		if key != "" {
			result[key] = value
		}
	}
	return result
}
