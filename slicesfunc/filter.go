package slicesfunc

// this helper function filters the given list using the value of the callback function

func Filter[I any](sliceList []I, f func(item I) bool) []I {
	var newList []I

	for _, sliceItem := range sliceList {
		if f(sliceItem) {
			newList = append(newList, sliceItem)
		}
	}

	return newList
}
