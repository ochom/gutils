package gttp

import "time"

// getTimeout returns the first provided timeout duration or a default of 10 seconds.
// Used internally by HTTP clients to determine request timeout.
func getTimeout(timeout ...time.Duration) time.Duration {
	if len(timeout) == 0 {
		return 10 * time.Second
	}

	return timeout[0]
}
