package helpers

import (
	"fmt"
	"net"
	"time"

	"github.com/ochom/gutils/logs"
)

// TestPorts is a list of ports that are used for testing
func GetPort(startingPort int) int {
	_, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", startingPort)), time.Second)
	if err == nil {
		logs.Warn("[🥵] port %d is not available trying another port: %d", startingPort, startingPort+1)
		return GetPort(startingPort + 1)
	}

	return startingPort
}
