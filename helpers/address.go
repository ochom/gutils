// Package helpers provides common utility functions for everyday programming tasks.
//
// This package includes utilities for:
//   - Network address handling
//   - CSV parsing and file reading
//   - Email validation
//   - Phone number parsing (Kenyan format)
//   - Password hashing and comparison
//   - OTP generation
//   - Conditional expressions
//
// Example usage:
//
//	// Find available port
//	addr := helpers.GetAvailableAddress(8080)
//	fmt.Println("Starting server on", addr)
//
//	// Hash a password
//	hashed := helpers.HashPassword("mysecretpassword")
//
//	// Verify password
//	if helpers.ComparePassword(hashed, "mysecretpassword") {
//		fmt.Println("Password matches!")
//	}
package helpers

import (
	"fmt"
	"net"
	"time"

	"github.com/ochom/gutils/logs"
)

// GetAvailableAddress finds the next available TCP port starting from the given port.
// If the specified port is in use, it recursively tries the next port until an available one is found.
//
// Returns the address in the format ":port" (e.g., ":8080").
//
// Example:
//
//	// Get an available port starting from 8080
//	addr := helpers.GetAvailableAddress(8080)
//	// addr = ":8080" if available, or ":8081", ":8082", etc.
//
//	// Use with HTTP server
//	http.ListenAndServe(helpers.GetAvailableAddress(3000), handler)
func GetAvailableAddress(port int) string {
	_, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), time.Second)
	if err == nil {
		logs.Warn("[🥵] address :%d is not available trying another port...", port)
		return GetAvailableAddress(port + 1)
	}

	return fmt.Sprintf(":%d", port)
}
