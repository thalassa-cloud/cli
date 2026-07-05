package shared

import (
	"context"
	"fmt"
	"time"
)

// WaitContext returns a child context capped by timeout for --wait flows.
func WaitContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc, error) {
	if timeout <= 0 {
		return nil, nil, fmt.Errorf("--wait-timeout must be greater than 0")
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	return ctx, cancel, nil
}
