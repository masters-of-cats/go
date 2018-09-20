// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin dragonfly freebsd netbsd openbsd

package unix

import (
	"syscall"
	"unsafe"
)

func Openat(fdat int, path string, flags int, perm uint32) (int, error) {
	var pathBytePointer *byte
	pathBytePointer, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, err
	}

	fdPointer, _, errNo := syscall.Syscall6(openatTrap, uintptr(fdat), uintptr(unsafe.Pointer(pathBytePointer)), uintptr(flags), uintptr(perm), 0, 0)
	fd := int(fdPointer)
	if errNo != 0 {
		return 0, errNo
	}

	return fd, nil
}
