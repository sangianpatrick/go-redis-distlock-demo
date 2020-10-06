package sync

import (
	"context"
	"time"
)

// Mutex is an interface for global mutual exclusion.
type Mutex interface {
	Lock(context.Context) (err error)
	Unlock(context.Context) (err error)
}

// DistributedLock is an interaface to instantiate the global mutual exclusion.
type DistributedLock interface {
	NewMutex(key string, maxRetries int, retryDelay, expiry time.Duration) (mutex Mutex)
}
