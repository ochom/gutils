package gttp

import "time"

func getTimeout(timeout ...time.Duration) time.Duration {
	if len(timeout) == 0 {
		return 10 * time.Second
	}

	return timeout[0]
}
