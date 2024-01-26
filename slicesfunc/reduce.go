package slicesfunc

// this helper function helps reduce a slice into whatever

func Reduce[T, M any](iterable []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range iterable {
		acc = f(acc, v)
	}
	return acc
}
