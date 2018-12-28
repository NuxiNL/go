// +build cloudabi

package os

import (
	"errors"
	"time"
)

type dirInfo struct{}

type file struct {
	name    string
	dirinfo *dirInfo // nil unless directory being read
}

func (f *File) readdir(n int) (fi []FileInfo, err error) {
	return nil, errors.New("XXX")
}

func (f *File) readdirnames(n int) (names []string, err error) {
	return nil, errors.New("XXX")
}

func NewFile(fd uintptr, name string) *File {
	return nil
}

func (f *File) checkValid(op string) error {
	if f == nil {
		return ErrInvalid
	}
	return nil
}

func (f *File) read(b []byte) (n int, err error) {
	return 0, errors.New("XXX")
}

func (f *File) write(b []byte) (n int, err error) {
	return 0, errors.New("XXX")
}

func (f *File) pread(b []byte, off int64) (n int, err error) {
	return 0, errors.New("XXX")
}

func (f *File) pwrite(b []byte, off int64) (n int, err error) {
	return 0, errors.New("XXX")
}

func (f *File) seek(offset int64, whence int) (ret int64, err error) {
	return 0, errors.New("XXX")
}

func (f *File) setDeadline(t time.Time) error {
	return errors.New("XXX")
}

// setReadDeadline sets the read deadline.
func (f *File) setReadDeadline(t time.Time) error {
	return errors.New("XXX")
}

// setWriteDeadline sets the write deadline.
func (f *File) setWriteDeadline(t time.Time) error {
	return errors.New("XXX")
}

func (f *File) chmod(mode FileMode) error {
	return errors.New("XXX")
}

func epipecheck(file *File, e error) {
}
