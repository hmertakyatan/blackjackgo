package services

import (
	"fmt"
	"sync/atomic"
)

type AtomicInt struct {
	value int32
}

func NewAtomicInt(value int32) *AtomicInt {
	return &AtomicInt{
		value: value,
	}
}

func (a *AtomicInt) ReadAtomicInt() string {
	return fmt.Sprintf("%d", a.value)
}

func (a *AtomicInt) SetAtomicInt(value int32) {
	atomic.StoreInt32(&a.value, value)
}

func (a *AtomicInt) GetAtomicInt() int32 {
	return atomic.LoadInt32(&a.value)
}

func (a *AtomicInt) IncrementAtomicIntValue() {
	currentValue := a.GetAtomicInt()
	a.SetAtomicInt(currentValue + 1)
}
