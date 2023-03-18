package main

import (
	"unsafe"

	"libdb.so/console"
)

func init() {
	if unsafe.Sizeof(uintptr(0)) != unsafe.Sizeof(uint32(0)) {
		panic("wasm is not 32-bit, what the fuck?")
	}
}

// read reads into the given bytes pointer that is of len size. It returns the
// number of bytes read. This should be used to read the terminal data.
//
//export read
func _read(ptr uint32, len uint32) int32 {
	n, err := terminal.Read(wasmBytes(ptr, len))
	if err != nil {
		return -1
	}

	return int32(n)
}

// write writes the given bytes pointer that is of len size into the stdin pipe.
// It returns false if the write failed. In most cases, the program should panic
// if that is the case.
//
// Note: the function MUST block until the write is complete. It also must not
// hold onto the bytes pointer after the write is complete.
//
//export write
func _write(ptr uint32, len uint32) bool {
	_, err := input.Write(wasmBytes(ptr, len))
	return err == nil
}

//export updateTerminal
func _updateTerminal(row, col, xpixel, ypixel uint32, sixel bool) {
	interp.UpdateTerminal(console.TerminalQuery{
		Width:  int(row),
		Height: int(col),
		XPixel: int(xpixel),
		YPixel: int(ypixel),
		SIXEL:  sixel,
	})
}

func wasmBytes(ptr, len uint32) []byte {
	// trust me bro.
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), len)
}

func wasmBytesPtr(b []byte) (ptr_, len_ uint32) {
	// trust me bro.
	return uint32(uintptr(unsafe.Pointer(&b[0]))), uint32(len(b))
}
