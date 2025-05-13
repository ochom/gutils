package uuid

import (
	"strconv"
	"strings"
)

// ShortID translates a number to a 5 character string
func ShortID(input int, size int) string {
	if size < 5 {
		size = 5
	}

	modulus := 1 << 16                        // Modulus for Feistel function
	rounds := 4                               // Number of Feistel rounds
	keys := []int{12345, 67890, 54321, 98765} // Round keys

	output := feistelNetwork(input, rounds, modulus, keys)

	res := strings.ToUpper(strconv.FormatInt(int64(output), 36))
	for len(res) < size {
		res += "Ø"
	}

	return res
}

// Feistel network for bijective mapping
func feistelNetwork(input, rounds, modulus int, keys []int) int {
	left := input >> 16     // Higher 16 bits
	right := input & 0xFFFF // Lower 16 bits

	for i := 0; i < rounds; i++ {
		newRight := left ^ feistelFunction(right, keys[i], modulus)
		left = right
		right = newRight
	}

	// Combine left and right parts
	return (right << 16) | left
}

// Feistel function: a simple example with modular arithmetic
func feistelFunction(data, key, modulus int) int {
	return (data*key ^ key) % modulus
}
