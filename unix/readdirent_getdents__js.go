// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js

package unix

// ReadDirent reads directory entries from fd and writes them into buf.
func ReadDirent(fd int, buf []byte) (n int, err error) {
	return Getdents(fd, buf)
}
