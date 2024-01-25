package slicesfunc

func Count(sliceList []string, item string) int {
	count := 0

	for _, s := range sliceList {
		if s == item {
			count++
		}
	}

	return count
}
