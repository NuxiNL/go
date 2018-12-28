// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cloudabi

package time

import (
	"errors"
)

func interrupt() {}

func open(name string) (uintptr, error) {
	return 0, errors.New("Cannot open files in this environment")
}

func read(fd uintptr, buf []byte) (int, error) {
	return 0, errors.New("Cannot read files in this environment")
}

func closefd(fd uintptr) {}

func preadn(fd uintptr, buf []byte, off int) error {
	return errors.New("Cannot read files in this environment")
}
