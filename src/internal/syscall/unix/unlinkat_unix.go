// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin dragonfly freebsd netbsd openbsd

package unix

import (
	"syscall"
	"unsafe"
)

func Unlinkat(dirfd int, path string, flags int) error {
	var pathBytePointer *byte
	pathBytePointer, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	_, _, errNo := syscall.Syscall(unlinkatTrap, uintptr(dirfd), uintptr(unsafe.Pointer(pathBytePointer)), uintptr(flags))
	if errNo != 0 {
		return errNo
	}

	return nil
}
