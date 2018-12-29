// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// CloudABI does not support command line arguments, so set argc to
// zero. Even though argv itself is of little use, use it to preserve
// the address of the auxiliary vector.
TEXT _rt0_amd64_cloudabi(SB),NOSPLIT,$-8
	MOVQ	DI, SI // argv == auxiliary vector
	MOVQ	$0, DI // argc == 0
	JMP	runtime·rt0_go(SB)

TEXT _rt0_amd64_cloudabi_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)
