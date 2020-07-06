// Copyright (C) 2020 Sachin Saini

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
)

// LruHandler implemnts methods to manage the lru cache
type LruHandler struct {
	cache   *Lru
	mem     *Manager
	memUsed int
}

// NewLruHandler instantiates a new lru handler
func NewLruHandler(m *Manager) *LruHandler {
	l := LruHandler{}
	l.cache = NewLru()
	l.mem = m

	return &l
}

// LRSet inserts an item to the lru
func (l *LruHandler) LRSet(key, value string) string {
	iter := 0
	for l.mem.CurrentlyUsed+len(value)+len(key) >= l.mem.MaxCapacity && iter < 100 {
		rkey, rvalue, err := l.cache.RemoveOldest()
		if err != nil {
			panic(err)
		}
		l.mem.Free(len(fmt.Sprint(rkey)) + len(fmt.Sprint(rvalue)))
		l.memUsed -= len(fmt.Sprint(rkey)) + len(fmt.Sprint(rvalue))
		iter++
	}
	if iter >= 100 {
		return MemoryOverLoad
	}
	l.mem.Allocate(len(key) + len(value))
	l.memUsed += len(key) + len(value)
	err := l.cache.Set(key, value)
	if err != nil {
		panic(err)
	}
	return Ok
}

// LRGet returns a key from the lru
func (l *LruHandler) LRGet(key string) string {
	value, err := l.cache.Get(key)
	if err != nil {
		return NoExist
	}
	return fmt.Sprint(value)
}

// LRRemove removes a key from the lru
func (l *LruHandler) LRRemove(key string) string {
	value, err := l.cache.Remove(key)
	if err != nil {
		return NoExist
	}
	l.mem.Free(len(fmt.Sprint(key)) + len(fmt.Sprint(value)))
	l.memUsed -= len(fmt.Sprint(key)) + len(fmt.Sprint(value))
	return fmt.Sprint(value)
}

// LRPurge clears the lru
func (l *LruHandler) LRPurge(key string) string {
	err := l.cache.Purge()
	l.mem.Free(l.memUsed)
	l.memUsed = 0
	if err != nil {
		panic(err)
	}
	return Ok
}