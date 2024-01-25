package slicesfunc

func Filter[I any](sliceList []I, f func(item I) bool) []I {
	var newList []I

	for _, sliceItem := range sliceList {
		if f(sliceItem) {
			newList = append(newList, sliceItem)
		}
	}

	return newList
}
