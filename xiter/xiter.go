package xiter

import "iter"

type FilterFunc[V any] func(V) bool
type FilterFunc2[K, V any] func(K, V) bool
type Map2Func[K, V, K2, V2 any] func(K, V) (K2, V2)
type FilterFunc2CanError[K, V any] func(K, V) (bool, error)

func Map2[K, V, K2, V2 any](s2 iter.Seq2[K, V], mapFunc Map2Func[K, V, K2, V2]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s2 {
			k2, v2 := mapFunc(k, v)
			if !yield(k2, v2) {
				return
			}
		}
	}
}

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
