package cast

func StringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Uint64Pointer(val uint64) *uint64 {
	return &val
}

func IntPointer(val int) *int {
	return &val
}
