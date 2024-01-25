package slicesfunc

func Any[I any](sliceList []I, f func(item I) bool) bool {
	for _, sliceItem := range sliceList {
		if f(sliceItem) {
			return true
		}
	}

	return false
}
