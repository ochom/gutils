package arrays

// Filter returns a new array with all elements that pass the test implemented by the provided function.
func Filter[S ~[]E, E any](items S, fn func(E) bool) []E {
	filteredItems := []E{}
	for _, value := range items {
		if fn(value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}

// Find returns the value of the first element in the provided array that satisfies the provided testing function.
func Find[S ~[]E, E any](items S, fn func(E) bool) E {
	var item E
	for _, value := range items {
		if fn(value) {
			return value
		}
	}
	return item
}

// Map returns a new array with the results of calling a provided function on every element in the provided array.
func Map[S, T any](items []S, fn func(S) T) []T {
	mappedItems := []T{}
	for _, value := range items {
		mappedItems = append(mappedItems, fn(value))
	}
	return mappedItems
}

// MapIndex returns a new array with the results of calling a provided function on every element in the provided array.
func MapIndex[S, T any](items []S, fn func(S, int) T) []T {
	mappedItems := []T{}
	for index, value := range items {
		mappedItems = append(mappedItems, fn(value, index))
	}
	return mappedItems
}

// Reduce applies a function against an accumulator and each element in the array (from left to right) to reduce it to a single value.
func Reduce[S, T any](items []S, fn func(T, S) T, acc T) T {
	for _, value := range items {
		acc = fn(acc, value)
	}
	return acc
}
