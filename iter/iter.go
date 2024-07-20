package iter

type (
	Seq[V any]     func(yield func(V) bool)
	Seq2[K, V any] func(yield func(K, V) bool)
)

type FilterFunc[V any] func(V) bool
type FilterFunc2[K, V any] func(K, V) bool
type FilterFunc2CanError[K, V any] func(K, V) (bool, error)

// Filter yields only values for which filterFunc returns true
func Filter[T any](s Seq[T], filterFunc FilterFunc[T]) Seq[T] {
	return func(yield func(T) bool) {
		s(func(value T) bool {
			if shouldYield := filterFunc(value); !shouldYield {
				return true
			}
			if yield(value) {
				return true
			}
			return false
		})
	}
}

func Filter2[K, V any](seq Seq2[K, V], filterFunc FilterFunc2[K, V]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(key K, value V) bool {
			if shouldYield := filterFunc(key, value); !shouldYield {
				return true
			}
			if yield(key, value) {
				return true
			}
			return false
		})
	}
}

func Filter2CanError[K, V any](seq Seq2[K, V], filterFunc FilterFunc2CanError[K, V]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(key K, value V) bool {
			shouldYield, err := filterFunc(key, value)
			if err != nil {
				return false
			}
			if !shouldYield {
				return true
			}
			if yield(key, value) {
				return true
			}
			return false
		})
	}
}
