// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build mips64 mips64le

package unix

const unlinkatTrap = uintptr(5253)
const openatTrap = uintptr(5247)

const AT_REMOVEDIR = 0x200
