package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

type LinearAllocator struct {
	data []byte
}

func NewLinearAllocator(capacity int) (LinearAllocator, error) {
	if capacity <= 0 {
		return LinearAllocator{}, errors.New("incorrect capacity")
	}

	return LinearAllocator{
		data: make([]byte, 0, capacity),
	}, nil
}

func (a *LinearAllocator) Allocate(size int) (unsafe.Pointer, error) {
	previousLength := len(a.data)
	newLength := previousLength + size

	if newLength > cap(a.data) {
		return nil, errors.New("not enough memory")
	}

	a.data = a.data[:newLength]
	pointer := unsafe.Pointer(&a.data[previousLength])
	return pointer, nil
}

// not supported by this kind of allocator
// func (a *LinearAllocator) Deallocate(pointer unsafe.Pointer) error {}

func (a *LinearAllocator) Free() {
	a.data = a.data[:0]
}

func store[T any](pointer unsafe.Pointer, value T) {
	*(*T)(pointer) = value
}

func load[T any](pointer unsafe.Pointer) T {
	return *(*T)(pointer)
}

func main() {
	const MB = 1 << 10 << 10
	allocator, err := NewLinearAllocator(MB)

	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialize allocator:", err)
		os.Exit(1)
	}
	defer allocator.Free()

	p1, _ := allocator.Allocate(2)
	p2, _ := allocator.Allocate(4)

	store[int16](p1, 100)
	store[int32](p2, 200)

	v1 := load[int16](p1)
	v2 := load[int32](p2)

	fmt.Println("value 1:", v1)
	fmt.Println("value 2:", v2)

	fmt.Println("address 1:", p1)
	fmt.Println("address 2:", p2)
}
