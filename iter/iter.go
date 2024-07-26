package iter

import "iter"

type FilterFunc[V any] func(V) bool
type FilterFunc2[K, V any] func(K, V) bool
type FilterFunc2CanError[K, V any] func(K, V) (bool, error)

func Filter[T any](seq iter.Seq[T], filterFunc FilterFunc[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range seq {
			if keep := filterFunc(item); keep {
				if !yield(item) {
					break
				}
			}
		}
	}
}

func Filter2[K, V any](seq iter.Seq2[K, V], filterFunc FilterFunc2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, value := range seq {
			if keep := filterFunc(key, value); keep {
				if !yield(key, value) {
					break
				}
			}
		}
	}
}

func Sample[K, V any](seq iter.Seq2[K, V], every int) iter.Seq2[K, V] {
	index := 0
	return func(yield func(K, V) bool) {
		for key, value := range seq {
			if index%every == 0 {
				if !yield(key, value) {
					break
				}
			}
			index++
		}
	}
}

func FilterWithError2[K, V any](seq iter.Seq2[K, V], filterFunc FilterFunc2CanError[K, V]) iter.Seq2[K, error] {
	return func(yield func(K, error) bool) {
		for key, value := range seq {
			keep, err := filterFunc(key, value)
			if err != nil {
				if !yield(key, err) {
					break
				}
			}

			if keep {
				if !yield(key, err) {
					break
				}
			}
		}
	}
}
