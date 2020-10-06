package sync

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
)

// RedsyncPool is an interface for building the mock of "redsync.Pool".
//
// Use this abstraction to auto-generate mock for unit testing instead of application usecase.
type RedsyncPool interface {
	Get(ctx context.Context) (redis.Conn, error)
}

// RedsyncConn is an interface for building the moc of "redsync.Conn".
//
// Use this abstraction to auto-generate mock for unit testing instead of application usecase.
type RedsyncConn interface {
	Get(name string) (string, error)
	Set(name string, value string) (bool, error)
	SetNX(name string, value string, expiry time.Duration) (bool, error)
	Eval(script *redis.Script, keysAndArgs ...interface{}) (interface{}, error)
	PTTL(name string) (time.Duration, error)
	Close() error
}

// RedsyncMutexAdapter is an adapter for redsync mutex concrete struct.
type RedsyncMutexAdapter struct {
	rmx *redsync.Mutex
}

// Lock returns error when it's fail to acquire the lock.
func (rsm *RedsyncMutexAdapter) Lock(ctx context.Context) (err error) {
	err = rsm.rmx.LockContext(ctx)
	return
}

// Unlock returns error when it's fail to release the lock.
func (rsm *RedsyncMutexAdapter) Unlock(ctx context.Context) (err error) {
	_, err = rsm.rmx.UnlockContext(ctx)
	return
}

// RedsyncAdapter is an adapter for redsync concrete struct.
type RedsyncAdapter struct {
	rs *redsync.Redsync
}

// NewRedsyncAdapter is a constructor to wrap redsync concrete struct to adapt as a general distributed lock.
func NewRedsyncAdapter(pools ...redis.Pool) DistributedLock {
	return &RedsyncAdapter{
		rs: redsync.New(pools...),
	}
}

// NewMutex returns mutex of redsync adapter
func (rsa *RedsyncAdapter) NewMutex(key string, maxRetries int, retryDelay, expiry time.Duration) (mutex Mutex) {
	redsyncMutex := rsa.rs.NewMutex(key,
		redsync.WithTries(maxRetries),
		redsync.WithRetryDelay(retryDelay),
		redsync.WithExpiry(expiry),
	)

	mutex = &RedsyncMutexAdapter{
		rmx: redsyncMutex,
	}

	return
}
