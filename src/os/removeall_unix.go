// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin

package os

import (
	"io"

	"internal/syscall/unix"
	"syscall"
)

func RemoveAll(path string) error {
	// Not allowed in unix
	if path == "" || endsWithDot(path) {
		return syscall.EINVAL
	}

	// RemoveAll recurses by deleting the path base from
	// its parent directory
	parentDir, base := splitPath(path)

	parent, err := Open(parentDir)
	if IsNotExist(err) {
		// If parent does not exist, base cannot exist. Fail silently
		return nil
	}
	if err != nil {
		return err
	}
	defer parent.Close()

	return removeAllFrom(parent, base)
}

func removeAllFrom(parentFile *File, path string) error {
	parentFd := int(parentFile.Fd())
	// Simple case: if Unlink (aka remove) works, we're done.
	err := unix.Unlinkat(parentFd, path, 0)
	if err == nil || IsNotExist(err) {
		// Already deleted by someone else
		return nil
	}

	// If not a "is directory" error, we have a problem
	if err != syscall.EISDIR && err != syscall.EPERM {
		return err
	}

	// Open the directory to recurse into
	fd, err := unix.Openat(parentFd, path, O_RDONLY, 0)
	if IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	file := NewFile(uintptr(fd), path)

	// Remove the directory's entries
	recurseErr := removeDirEntries(file)

	file.Close()

	// Remove the directory itself
	unlinkError := unix.Unlinkat(parentFd, path, unix.AT_REMOVEDIR)
	if unlinkError == nil || IsNotExist(unlinkError) {
		return nil
	}

	if recurseErr != nil {
		return recurseErr
	}
	return unlinkError
}

func removeDirEntries(file *File) error {
	var removeErr error
	for {
		const request = 1024
		names, readErr := file.Readdirnames(request)
		// Errors other than EOF should stop us from continuing
		if readErr != nil && readErr != io.EOF {
			return readErr
		}

		for _, name := range names {
			err := removeAllFrom(file, name)
			if err != nil {
				removeErr = err
			}
		}
		// Return the most recently occured RemoveAll error when
		// the end of the directory is reached
		if len(names) < request {
			return removeErr
		}
	}
}

func endsWithDot(path string) bool {
	last := len(path) - 1
	if path[last] == '.' {
		return true
	}

	return false
}
