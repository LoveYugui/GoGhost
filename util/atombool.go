package util

import (
	"fmt"
	"sync"
	"sync/atomic"
	_ "github.com/GoGhost/constants"
)

type AtomicInt64 int64

func NewAtomicInt64(initialValue int64) *AtomicInt64 {
	a := AtomicInt64(initialValue)
	return &a
}

func (a *AtomicInt64) Get() int64 {
	return int64(*a)
}

func (a *AtomicInt64) Set(newValue int64) {
	atomic.StoreInt64((*int64)(a), newValue)
}

func (a *AtomicInt64) GetAndSet(newValue int64) int64 {
	for {
		current := a.Get()
		if a.CompareAndSet(current, newValue) {
			return current
		}
	}
}

func (a *AtomicInt64) CompareAndSet(expect, update int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(a), expect, update)
}

func (a *AtomicInt64) GetAndIncrement() int64 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}

}

func (a *AtomicInt64) GetAndDecrement() int64 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *AtomicInt64) GetAndAdd(delta int64) int64 {
	for {
		current := a.Get()
		next := current + delta
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *AtomicInt64) IncrementAndGet() int64 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *AtomicInt64) DecrementAndGet() int64 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *AtomicInt64) AddAndGet(delta int64) int64 {
	for {
		current := a.Get()
		next := current + delta
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *AtomicInt64) String() string {
	return fmt.Sprintf("%d", a.Get())
}

type AtomicInt32 int32

func NewAtomicInt32(initialValue int32) *AtomicInt32 {
	a := AtomicInt32(initialValue)
	return &a
}

func (a *AtomicInt32) Get() int32 {
	return int32(*a)
}

func (a *AtomicInt32) Set(newValue int32) {
	atomic.StoreInt32((*int32)(a), newValue)
}

func (a *AtomicInt32) GetAndSet(newValue int32) (oldValue int32) {
	for {
		oldValue = a.Get()
		if a.CompareAndSet(oldValue, newValue) {
			return
		}
	}
}

func (a *AtomicInt32) CompareAndSet(expect, update int32) bool {
	return atomic.CompareAndSwapInt32((*int32)(a), expect, update)
}

func (a *AtomicInt32) GetAndIncrement() int32 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}

}

func (a *AtomicInt32) GetAndDecrement() int32 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *AtomicInt32) GetAndAdd(delta int32) int32 {
	for {
		current := a.Get()
		next := current + delta
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *AtomicInt32) IncrementAndGet() int32 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *AtomicInt32) DecrementAndGet() int32 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *AtomicInt32) AddAndGet(delta int32) int32 {
	for {
		current := a.Get()
		next := current + delta
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *AtomicInt32) String() string {
	return fmt.Sprintf("%d", a.Get())
}

type AtomicBoolean int32

func NewAtomicBoolean(initialValue bool) *AtomicBoolean {
	var a AtomicBoolean
	if initialValue {
		a = AtomicBoolean(1)
	} else {
		a = AtomicBoolean(0)
	}
	return &a
}

func (a *AtomicBoolean) Get() bool {
	return atomic.LoadInt32((*int32)(a)) != 0
}

func (a *AtomicBoolean) Set(newValue bool) {
	if newValue {
		atomic.StoreInt32((*int32)(a), 1)
	} else {
		atomic.StoreInt32((*int32)(a), 0)
	}
}

func (a *AtomicBoolean) CompareAndSet(oldValue, newValue bool) bool {
	var o int32
	var n int32
	if oldValue {
		o = 1
	} else {
		o = 0
	}
	if newValue {
		n = 1
	} else {
		n = 0
	}
	return atomic.CompareAndSwapInt32((*int32)(a), o, n)
}

func (a *AtomicBoolean) GetAndSet(newValue bool) bool {
	for {
		current := a.Get()
		if a.CompareAndSet(current, newValue) {
			return current
		}
	}
}

func (a *AtomicBoolean) String() string {
	return fmt.Sprintf("%t", a.Get())
}

type Item struct {
	Key, Value interface{}
}

type SyncMap struct {
	shard map[interface{}]interface{}
	sync.RWMutex
}

func NewSyncMap()*SyncMap {
	return &SyncMap{
		shard: make(map[interface{}]interface{}, 32),
	}
}

func (sm *SyncMap) Put(k, v interface{}) {
	sm.Lock()
	sm.shard[k] = v
	sm.Unlock()
}

func (sm *SyncMap) Get(k interface{}) (interface{}, bool) {
	sm.RLock()
	v, ok := sm.shard[k]
	sm.RUnlock()
	return v, ok
}

func (sm *SyncMap) Size() int {
	var ret int
	sm.RLock()
	ret = len(sm.shard)
	sm.RUnlock()
	return ret
}

func (sm *SyncMap) Remove(k interface{}) {
	sm.Lock()
	delete(sm.shard, k)
	sm.Unlock()
}

func (sm *SyncMap) Clear() {
	sm.Lock()
	sm.shard = make(map[interface{}]interface{})
	sm.Unlock()
}

//获取一个副本，尽量少用
func (sm *SyncMap) All() map[interface{}]interface{} {
	sm.Lock()
	newMap := make(map[interface{}]interface{})
	for k, v := range sm.shard {
		newMap[k] = DeepCopy(v)
	}

	sm.Unlock()

	return newMap
}

//
func DeepCopy(value interface{}) interface{} {
	//非map[string]interface{}的会有问题
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}