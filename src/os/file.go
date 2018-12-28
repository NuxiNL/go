// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package os provides a platform-independent interface to operating system
// functionality. The design is Unix-like, although the error handling is
// Go-like; failing calls return values of type error rather than error numbers.
// Often, more information is available within the error. For example,
// if a call that takes a file name fails, such as Open or Stat, the error
// will include the failing file name when printed and will be of type
// *PathError, which may be unpacked for more information.
//
// The os interface is intended to be uniform across all operating systems.
// Features not generally available appear in the system-specific package syscall.
//
// Here is a simple example, opening a file and reading some of it.
//
//	file, err := os.Open("file.go") // For read access.
//	if err != nil {
//		log.Fatal(err)
//	}
//
// If the open fails, the error string will be self-explanatory, like
//
//	open file.go: no such file or directory
//
// The file's data can then be read into a slice of bytes. Read and
// Write take their byte counts from the length of the argument slice.
//
//	data := make([]byte, 100)
//	count, err := file.Read(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("read %d bytes: %q\n", count, data[:count])
//
package os

import (
	"errors"
	"internal/poll"
	"io"
	"syscall"
	"time"
)

// Name returns the name of the file as presented to Open.
func (f *File) Name() string { return f.name }

// Flags to OpenFile wrapping those of the underlying system. Not all
// flags may be implemented on a given system.
const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
)

// Seek whence values.
//
// Deprecated: Use io.SeekStart, io.SeekCurrent, and io.SeekEnd.
const (
	SEEK_SET int = 0 // seek relative to the origin of the file
	SEEK_CUR int = 1 // seek relative to the current offset
	SEEK_END int = 2 // seek relative to the end
)

// LinkError records an error during a link or symlink or rename
// system call and the paths that caused it.
type LinkError struct {
	Op  string
	Old string
	New string
	Err error
}

func (e *LinkError) Error() string {
	return e.Op + " " + e.Old + " " + e.New + ": " + e.Err.Error()
}

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (f *File) Read(b []byte) (n int, err error) {
	if err := f.checkValid("read"); err != nil {
		return 0, err
	}
	n, e := f.read(b)
	return n, f.wrapErr("read", e)
}

// ReadAt reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b).
// At end of file, that error is io.EOF.
func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	if err := f.checkValid("read"); err != nil {
		return 0, err
	}

	if off < 0 {
		return 0, &PathError{"readat", f.name, errors.New("negative offset")}
	}

	for len(b) > 0 {
		m, e := f.pread(b, off)
		if e != nil {
			err = f.wrapErr("read", e)
			break
		}
		n += m
		b = b[m:]
		off += int64(m)
	}
	return
}

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (f *File) Write(b []byte) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	n, e := f.write(b)
	if n < 0 {
		n = 0
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	epipecheck(f, e)

	if e != nil {
		err = f.wrapErr("write", e)
	}

	return n, err
}

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}

	if off < 0 {
		return 0, &PathError{"writeat", f.name, errors.New("negative offset")}
	}

	for len(b) > 0 {
		m, e := f.pwrite(b, off)
		if e != nil {
			err = f.wrapErr("write", e)
			break
		}
		n += m
		b = b[m:]
		off += int64(m)
	}
	return
}

// Seek sets the offset for the next Read or Write on file to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error, if any.
// The behavior of Seek on a file opened with O_APPEND is not specified.
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	if err := f.checkValid("seek"); err != nil {
		return 0, err
	}
	r, e := f.seek(offset, whence)
	if e == nil && f.dirinfo != nil && r != 0 {
		e = syscall.EISDIR
	}
	if e != nil {
		return 0, f.wrapErr("seek", e)
	}
	return r, nil
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (f *File) WriteString(s string) (n int, err error) {
	return f.Write([]byte(s))
}

// wrapErr wraps an error that occurred during an operation on an open file.
// It passes io.EOF through unchanged, otherwise converts
// poll.ErrFileClosing to ErrClosed and wraps the error in a PathError.
func (f *File) wrapErr(op string, err error) error {
	if err == nil || err == io.EOF {
		return err
	}
	if err == poll.ErrFileClosing {
		err = ErrClosed
	}
	return &PathError{op, f.name, err}
}

// Chmod changes the mode of the file to mode.
// If there is an error, it will be of type *PathError.
func (f *File) Chmod(mode FileMode) error { return f.chmod(mode) }

// SetDeadline sets the read and write deadlines for a File.
// It is equivalent to calling both SetReadDeadline and SetWriteDeadline.
//
// Only some kinds of files support setting a deadline. Calls to SetDeadline
// for files that do not support deadlines will return ErrNoDeadline.
// On most systems ordinary files do not support deadlines, but pipes do.
//
// A deadline is an absolute time after which I/O operations fail with an
// error instead of blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or Write.
// After a deadline has been exceeded, the connection can be refreshed
// by setting a deadline in the future.
//
// An error returned after a timeout fails will implement the
// Timeout method, and calling the Timeout method will return true.
// The PathError and SyscallError types implement the Timeout method.
// In general, call IsTimeout to test whether an error indicates a timeout.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (f *File) SetDeadline(t time.Time) error {
	return f.setDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls and any
// currently-blocked Read call.
// A zero value for t means Read will not time out.
// Not all files support setting deadlines; see SetDeadline.
func (f *File) SetReadDeadline(t time.Time) error {
	return f.setReadDeadline(t)
}

// SetWriteDeadline sets the deadline for any future Write calls and any
// currently-blocked Write call.
// Even if Write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
// Not all files support setting deadlines; see SetDeadline.
func (f *File) SetWriteDeadline(t time.Time) error {
	return f.setWriteDeadline(t)
}
