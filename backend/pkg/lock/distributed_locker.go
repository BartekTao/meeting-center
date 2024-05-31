package lock

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type DistributedLocker interface {
	Lock(key string) (*redsync.Mutex, error)
	Unlock(mutex *redsync.Mutex) (bool, error)
	TryLockWithWait(key string, retryInterval time.Duration, maxRetries int) (*redsync.Mutex, error)
}

type RedsyncLocker struct {
	rs *redsync.Redsync
}

func NewRedsyncLocker(rs *redsync.Redsync) *RedsyncLocker {
	return &RedsyncLocker{
		rs: rs,
	}
}

func (r *RedsyncLocker) TryLockWithWait(key string, retryInterval time.Duration, maxRetries int) (*redsync.Mutex, error) {
	for attempts := 0; attempts < maxRetries; attempts++ {
		mutex := r.rs.NewMutex(key)
		if err := mutex.Lock(); err == nil {
			return mutex, nil
		}
		time.Sleep(retryInterval)
	}
	return nil, fmt.Errorf("failed to acquire lock after %d attempts", maxRetries)
}

func (r *RedsyncLocker) Lock(key string) (*redsync.Mutex, error) {
	mutex := r.rs.NewMutex(key)
	err := mutex.Lock()
	if err != nil {
		return nil, err
	}
	return mutex, nil
}

func (r *RedsyncLocker) Unlock(mutex *redsync.Mutex) (bool, error) {
	return mutex.Unlock()
}
