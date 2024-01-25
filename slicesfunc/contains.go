package slicesfunc

func Contains[T comparable](value T, sliceList []T) bool {
	for _, v := range sliceList {
		if v == value {
			return true
		}
	}

	return false
}
