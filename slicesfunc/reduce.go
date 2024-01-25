package slicesfunc

func Reduce[T, M any](iterable []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range iterable {
		acc = f(acc, v)
	}
	return acc
}
