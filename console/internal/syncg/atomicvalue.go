package syncg

import "sync/atomic"

// AtomicValue is a generic wrapper around atomic.Value.
type AtomicValue[T any] atomic.Value

// Load returns the value stored in the AtomicValue.
func (v *AtomicValue[T]) Load() (T, bool) {
	val, ok := (*atomic.Value)(v).Load().(T)
	if ok {
		return val, true
	}
	var z T
	return z, false
}

// Store stores the given value in the AtomicValue.
func (v *AtomicValue[T]) Store(val T) {
	(*atomic.Value)(v).Store(val)
}

// Swap swaps the value stored in the AtomicValue with the given value and
// returns the old value.
func (v *AtomicValue[T]) Swap(new T) T {
	return (*atomic.Value)(v).Swap(new).(T)
}

// CompareAndSwap compares the value stored in the AtomicValue with the given
// value and, if they are equal, stores the new value and returns true. If they
// are not equal, it returns false.
func (v *AtomicValue[T]) CompareAndSwap(old, new T) bool {
	return (*atomic.Value)(v).CompareAndSwap(old, new)
}
