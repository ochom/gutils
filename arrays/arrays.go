// Package arrays provides functional programming utilities for slice manipulation in Go.
//
// This package implements common array/slice operations similar to JavaScript array methods,
// making it easier for developers familiar with JavaScript to work with Go slices.
// All functions are generic and work with any slice type.
//
// Note: While these utilities are convenient, idiomatic Go often prefers explicit loops
// for performance-critical code.
//
// Example usage:
//
//	numbers := []int{1, 2, 3, 4, 5}
//
//	// Filter even numbers
//	evens := arrays.Filter(numbers, func(n int) bool { return n%2 == 0 })
//	// evens = [2, 4]
//
//	// Map to squares
//	squares := arrays.Map(numbers, func(n int) int { return n * n })
//	// squares = [1, 4, 9, 16, 25]
//
//	// Sum all numbers
//	sum := arrays.Reduce(numbers, 0, func(n, acc int) int { return acc + n })
//	// sum = 15
package arrays

// Filter returns a new slice containing only the elements for which the predicate function returns true.
//
// The function iterates through each element in the input slice and includes it in the result
// only if fn(element) returns true.
//
// Example:
//
//	// Filter positive numbers
//	numbers := []int{-2, -1, 0, 1, 2}
//	positives := arrays.Filter(numbers, func(n int) bool { return n > 0 })
//	// positives = [1, 2]
//
//	// Filter users by age
//	type User struct { Name string; Age int }
//	users := []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 17}}
//	adults := arrays.Filter(users, func(u User) bool { return u.Age >= 18 })
//	// adults = [{Name: "Alice", Age: 25}]
func Filter[S ~[]E, E any](items S, fn func(E) bool) []E {
	filteredItems := []E{}
	for _, value := range items {
		if fn(value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}

// Find returns the first element in the slice that satisfies the provided predicate function.
// If no element satisfies the predicate, the zero value of type E is returned.
//
// Example:
//
//	// Find first even number
//	numbers := []int{1, 3, 4, 5, 6}
//	firstEven := arrays.Find(numbers, func(n int) bool { return n%2 == 0 })
//	// firstEven = 4
//
//	// Find user by name
//	type User struct { Name string; Age int }
//	users := []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}}
//	user := arrays.Find(users, func(u User) bool { return u.Name == "Bob" })
//	// user = {Name: "Bob", Age: 30}
func Find[S ~[]E, E any](items S, fn func(E) bool) E {
	var item E
	for _, value := range items {
		if fn(value) {
			return value
		}
	}
	return item
}

// Map transforms each element in the input slice using the provided transformation function
// and returns a new slice containing the transformed values.
//
// Example:
//
//	// Convert integers to strings
//	numbers := []int{1, 2, 3}
//	strings := arrays.Map(numbers, func(n int) string { return fmt.Sprintf("%d", n) })
//	// strings = ["1", "2", "3"]
//
//	// Extract names from structs
//	type User struct { Name string; Age int }
//	users := []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}}
//	names := arrays.Map(users, func(u User) string { return u.Name })
//	// names = ["Alice", "Bob"]
func Map[S, T any](items []S, fn func(S) T) []T {
	mappedItems := []T{}
	for _, value := range items {
		mappedItems = append(mappedItems, fn(value))
	}
	return mappedItems
}

// MapIndex is similar to Map but also passes the element index to the transformation function.
// This is useful when the transformation depends on the position of the element.
//
// Example:
//
//	// Create indexed items
//	items := []string{"a", "b", "c"}
//	result := arrays.MapIndex(items, func(s string, i int) string {
//		return fmt.Sprintf("%d: %s", i, s)
//	})
//	// result = ["0: a", "1: b", "2: c"]
func MapIndex[S, T any](items []S, fn func(S, int) T) []T {
	mappedItems := []T{}
	for index, value := range items {
		mappedItems = append(mappedItems, fn(value, index))
	}
	return mappedItems
}

// Reduce applies a reducer function to each element of the slice, accumulating the results
// into a single value. The accumulator (acc) is the initial value and carries the
// intermediate result through each iteration.
//
// Example:
//
//	// Sum all numbers
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := arrays.Reduce(numbers, 0, func(n, acc int) int { return acc + n })
//	// sum = 15
//
//	// Concatenate strings
//	words := []string{"Hello", " ", "World"}
//	result := arrays.Reduce(words, "", func(s, acc string) string { return acc + s })
//	// result = "Hello World"
//
//	// Build a map from slice
//	type User struct { ID string; Name string }
//	users := []User{{ID: "1", Name: "Alice"}, {ID: "2", Name: "Bob"}}
//	userMap := arrays.Reduce(users, map[string]string{}, func(u User, acc map[string]string) map[string]string {
//		acc[u.ID] = u.Name
//		return acc
//	})
//	// userMap = {"1": "Alice", "2": "Bob"}
func Reduce[S, T any](items []S, acc T, fn func(S, T) T) T {
	for _, value := range items {
		acc = fn(value, acc)
	}
	return acc
}

// ForEach executes the provided function once for each element in the slice.
// Unlike Map, it does not return a new slice and is used for side effects.
//
// Example:
//
//	// Print each element
//	numbers := []int{1, 2, 3}
//	arrays.ForEach(numbers, func(n int) {
//		fmt.Println(n)
//	})
//
//	// Process each user
//	type User struct { Name string; Email string }
//	users := []User{{Name: "Alice", Email: "alice@example.com"}}
//	arrays.ForEach(users, func(u User) {
//		sendWelcomeEmail(u.Email)
//	})
func ForEach[S any](items []S, fn func(S)) {
	for _, value := range items {
		fn(value)
	}
}
