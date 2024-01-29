package slicesfunc

// this helper function checks if the value exists in the given list

func Contains[T comparable](value T, sliceList []T) bool {
	for _, v := range sliceList {
		if v == value {
			return true
		}
	}

	return false
}
