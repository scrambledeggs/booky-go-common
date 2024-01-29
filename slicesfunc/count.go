package slicesfunc

// this helper function counts the occurence of the given element in the given list

func Count(sliceList []string, item string) int {
	count := 0

	for _, s := range sliceList {
		if s == item {
			count++
		}
	}

	return count
}
