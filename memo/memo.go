package memo

// hashKeyFunc represents the type of a function that generates a key for caching based on the function's input.
type hashKeyFunc[K any] func(K) string

type CloseFunc func()

// Memoize takes a function and a hash key function. It returns a memoized version of the function.
func Memoize[K any, V any](f func(K) V, hashKey hashKeyFunc[K]) (func(K) V, CloseFunc) {
	cache := make(map[string]V)
	return func(input K) V {
			hashedKey := hashKey(input)
			if val, found := cache[hashedKey]; found {
				return val
			}
			val := f(input)
			cache[hashedKey] = val
			return val
		}, func() {
			for _, value := range cache {
				if closer, ok := any(value).(interface{ Close() error }); ok {
					closer.Close()
				}
			}
		}
}
