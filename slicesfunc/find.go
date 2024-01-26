package slicesfunc

// this helper function returns the element if it exists in the list. otherwise, return nil record. the second element of the return indicates if the first element is usable

func Find[I any](sliceList []I, f func(item I) bool) (*I, bool) {
	for _, sliceItem := range sliceList {
		if f(sliceItem) {
			return &sliceItem, true
		}
	}

	return nil, false
}
