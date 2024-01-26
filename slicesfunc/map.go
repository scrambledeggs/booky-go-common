package slicesfunc

// this helper function helps modify all elements in the slice

func Map[I any, R any](sliceList []I, f func(item I) R) []R {
	var newList []R

	for _, sliceItem := range sliceList {
		newList = append(newList, f(sliceItem))
	}

	return newList
}
