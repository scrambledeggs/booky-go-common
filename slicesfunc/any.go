package slicesfunc

// this helper function checks if the value of the callback function exists in the given list

func Any[I any](sliceList []I, f func(item I) bool) bool {
	for _, sliceItem := range sliceList {
		if f(sliceItem) {
			return true
		}
	}

	return false
}
