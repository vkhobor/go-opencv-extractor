package uiter

import (
	i "iter"
)

func Tap[K, V any](seq i.Seq2[K, V], fn func(K, V)) i.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for value, err := range seq {
			fn(value, err)
			if !yield(value, err) {
				return
			}
		}
	}
}
