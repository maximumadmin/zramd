package utsname

func parseCharSlice(data []uint8) string {
	b := make([]byte, 0, len(data))
	for _, v := range data {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}
