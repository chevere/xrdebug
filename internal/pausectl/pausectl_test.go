/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package pausectl

import (
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	manager := NewManager(5*time.Minute, 1*time.Minute)
	lockID := "test-lock"
	t.Run("create new lock", func(t *testing.T) {
		lock, err := manager.New(lockID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if lock.ID != lockID {
			t.Errorf("Expected ID %s, got %s", lockID, lock.ID)
		}
		if lock.Stop {
			t.Error("Expected Stop to be false")
		}
	})
	t.Run("create duplicate lock", func(t *testing.T) {
		_, err := manager.New(lockID)
		if err != ErrLockExists {
			t.Errorf("Expected ErrLockExists, got %v", err)
		}
	})
	t.Run("get existing lock", func(t *testing.T) {
		lock, err := manager.Get(lockID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if lock.ID != lockID {
			t.Errorf("Expected ID %s, got %s", lockID, lock.ID)
		}
	})
	t.Run("get non-existent lock", func(t *testing.T) {
		_, err := manager.Get("non-existent")
		if err != ErrLockNotFound {
			t.Errorf("Expected ErrLockNotFound, got %v", err)
		}
	})
	t.Run("update lock", func(t *testing.T) {
		lock, err := manager.Update(lockID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !lock.Stop {
			t.Error("Expected Stop to be true")
		}
	})
	t.Run("delete lock", func(t *testing.T) {
		manager.Delete(lockID)
		_, err := manager.Get(lockID)
		if err != ErrLockNotFound {
			t.Errorf("Expected ErrLockNotFound, got %v", err)
		}
	})
	t.Run("update non-existent lock", func(t *testing.T) {
		_, err := manager.Update("non-existent")
		if err != ErrLockNotFound {
			t.Errorf("Expected ErrLockNotFound, got %v", err)
		}
	})
}
