package slicesfunc

func Find[I any](sliceList []I, f func(item I) bool) (*I, bool) {
	for _, sliceItem := range sliceList {
		if f(sliceItem) {
			return &sliceItem, true
		}
	}

	return nil, false
}
