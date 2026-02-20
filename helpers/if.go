package helpers

// If returns trueValue if condition is true, otherwise returns falseValue
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}
