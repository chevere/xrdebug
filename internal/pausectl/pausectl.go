/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

// Package pausectl provides functionality to manage pause locks with expiration.
package pausectl

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	ErrLockExists   = errors.New("lock already exists")
	ErrLockNotFound = errors.New("lock not found")
)

// Lock represents a pause lock with an identifier and stop status
type Lock struct {
	ID   string `json:"-"`
	Stop bool   `json:"stop"`
}

// Manager handles the creation and management of pause locks
type Manager struct {
	cache      *cache.Cache
	expiration time.Duration
}

// NewManager creates a new Manager with the specified expiration and cleanup intervals
func NewManager(expiration, cleanupInterval time.Duration) *Manager {
	return &Manager{
		cache:      cache.New(expiration, cleanupInterval),
		expiration: expiration,
	}
}

// New creates a new Lock with the specified ID
func (m *Manager) New(id string) (*Lock, error) {
	if _, found := m.cache.Get(id); found {
		return nil, ErrLockExists
	}
	m.cache.Set(id, false, m.expiration)
	return &Lock{id, false}, nil
}

// Get retrieves an existing Lock by its ID
func (m *Manager) Get(id string) (*Lock, error) {
	value, found := m.cache.Get(id)
	if !found {
		return nil, ErrLockNotFound
	}
	return &Lock{id, value.(bool)}, nil
}

// Update sets the stop status of a Lock to true
func (m *Manager) Update(id string) (*Lock, error) {
	if _, found := m.cache.Get(id); !found {
		return nil, ErrLockNotFound
	}
	m.cache.Set(id, true, m.expiration)
	return &Lock{id, true}, nil
}

// Delete removes a Lock from the manager
func (m *Manager) Delete(id string) {
	m.cache.Delete(id)
}
