package helpers

// If implements a ternary-like conditional expression.
// Returns trueValue if condition is true, otherwise returns falseValue.
//
// This is useful for inline conditional assignments that would otherwise
// require an if-else statement.
//
// Example:
//
//	// Simple conditional
//	status := helpers.If(user.Active, "active", "inactive")
//
//	// With function calls (note: both values are evaluated)
//	message := helpers.If(count > 0, "Found items", "No items")
//
//	// Numeric values
//	max := helpers.If(a > b, a, b)
//
//	// Struct values
//	config := helpers.If(isProd, prodConfig, devConfig)
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}
