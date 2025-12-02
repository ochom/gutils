package helpers

import (
	"fmt"
	"net"
	"time"

	"github.com/ochom/gutils/logs"
)

// GetAvailableAddress returns the next available address e.g :8080
func GetAvailableAddress(port int) string {
	_, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), time.Second)
	if err == nil {
		logs.Warn("[🥵] address :%d is not available trying another port...", port)
		return GetAvailableAddress(port + 1)
	}

	return fmt.Sprintf(":%d", port)
}
