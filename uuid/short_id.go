package uuid

import (
	"strconv"
	"strings"
)

// ShortID converts an integer to a unique, short alphanumeric string.
// Uses a Feistel network to create a bijective (reversible) mapping,
// ensuring each input produces a unique output and vice versa.
//
// The output is a base-36 uppercase string (0-9, A-Z), making it
// URL-safe and human-readable.
//
// Example:
//
//	// Generate short IDs from sequential numbers
//	id1 := uuid.ShortID(1)     // e.g., "ABC12"
//	id2 := uuid.ShortID(2)     // e.g., "XYZ89"
//	id3 := uuid.ShortID(12345) // e.g., "P3QR7"
//
//	// Use with auto-incrementing database IDs
//	publicID := uuid.ShortID(user.InternalID)
//	// Share publicID instead of revealing internal sequence
//
//	// URL-friendly order numbers
//	orderCode := uuid.ShortID(order.ID)
//	fmt.Sprintf("https://shop.com/order/%s", orderCode)
func ShortID(input int) string {
	modulus := 1 << 16                        // Modulus for Feistel function
	rounds := 4                               // Number of Feistel rounds
	keys := []int{12345, 67890, 54321, 98765} // Round keys

	output := feistelNetwork(input, rounds, modulus, keys)

	res := strings.ToUpper(strconv.FormatInt(int64(output), 36))
	return res
}

// feistelNetwork implements a Feistel cipher for bijective integer mapping.
// Ensures each input maps to exactly one output and the mapping is reversible.
func feistelNetwork(input, rounds, modulus int, keys []int) int {
	left := input >> 16     // Higher 16 bits
	right := input & 0xFFFF // Lower 16 bits

	for i := range rounds {
		newRight := left ^ feistelFunction(right, keys[i], modulus)
		left = right
		right = newRight
	}

	// Combine left and right parts
	return (right << 16) | left
}

// feistelFunction is the round function for the Feistel network.
func feistelFunction(data, key, modulus int) int {
	return (data*key ^ key) % modulus
}
