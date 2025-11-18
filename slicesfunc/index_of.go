package slicesfunc

func IndexOf[I any, R int](slice []I, f func(item I) bool) int {
	for i, sliceItem := range slice {
		if f(sliceItem) {
			return i
		}
	}
	return -1
}
