// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cloudabi

package os

import "errors"

func executable() (string, error) {
	return "", errors.New("executable path is unknown in this environment")
}
