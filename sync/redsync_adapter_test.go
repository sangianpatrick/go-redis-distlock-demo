package sync_test

import (
	"context"
	"testing"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"

	"github.com/sangianpatrick/go-redis-distlock-demo/sync"
	"github.com/sangianpatrick/go-redis-distlock-demo/sync/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRedsyncAdapter(t *testing.T) {
	pool := &mocks.RedsyncPool{}
	distlock := sync.NewRedsyncAdapter(pool)

	assert.NotNil(t, distlock)
}

func TestNewMutex(t *testing.T) {
	testKey := "test-key"
	pool := &mocks.RedsyncPool{}
	distlock := sync.NewRedsyncAdapter(pool)

	mutex := distlock.NewMutex(testKey, 5, time.Millisecond*200, time.Second*2)

	assert.NotNil(t, mutex)
}

func TestLock(t *testing.T) {
	conn := &mocks.RedsyncConn{}
	conn.On("SetNX", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration")).Return(true, nil).Once()
	conn.On("Close").Return(nil).Once()

	pool := &mocks.RedsyncPool{}
	pool.On("Get", mock.Anything).Return(conn, nil)

	distlock := sync.NewRedsyncAdapter(pool)

	testKey := "test-key"
	mutex := distlock.NewMutex(testKey, 5, time.Millisecond*200, time.Second*1)
	mutex.Lock(context.TODO())

	conn.AssertExpectations(t)
	pool.AssertExpectations(t)
}

func TestUnlock(t *testing.T) {
	rCmd := redisv8.NewCmd(context.TODO())

	conn := &mocks.RedsyncConn{}
	conn.On("SetNX", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration")).Return(true, nil).Once()
	conn.On("Eval", mock.AnythingOfType("*redis.Script"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(rCmd, nil).Once()
	conn.On("Close").Return(nil)

	pool := &mocks.RedsyncPool{}
	pool.On("Get", mock.Anything).Return(conn, nil)

	distlock := sync.NewRedsyncAdapter(pool)

	testKey := "test-key"
	mutex := distlock.NewMutex(testKey, 5, time.Millisecond*200, time.Second*1)

	mutex.Lock(context.TODO())
	mutex.Unlock(context.TODO())

	conn.AssertExpectations(t)
	pool.AssertExpectations(t)
}
