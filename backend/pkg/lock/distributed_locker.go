package lock

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type DistributedLocker interface {
	Lock(key string) (bool, error)
	Unlock(key string) error
	TryLockWithWait(key string, retryInterval time.Duration, maxRetries int) (bool, error)
}

type RedsyncLocker struct {
	rs *redsync.Redsync
}

func NewRedsyncLocker(rs *redsync.Redsync) *RedsyncLocker {
	return &RedsyncLocker{
		rs: rs,
	}
}

func (r *RedsyncLocker) TryLockWithWait(key string, retryInterval time.Duration, maxRetries int) (bool, error) {
	for attempts := 0; attempts < maxRetries; attempts++ {
		mutex := r.rs.NewMutex(key)
		if err := mutex.Lock(); err == nil {
			return true, nil
		}
		time.Sleep(retryInterval)
	}
	return false, fmt.Errorf("failed to acquire lock after %d attempts", maxRetries)
}

func (r *RedsyncLocker) Lock(key string) (bool, error) {
	mutex := r.rs.NewMutex(key)
	err := mutex.Lock()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedsyncLocker) Unlock(key string) error {
	mutex := r.rs.NewMutex(key)
	_, err := mutex.Unlock()
	return err
}
