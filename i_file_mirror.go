package gofilemirror

import (
	"io"
	"os"
)

// implements IFile
type IFileMirror interface {
	CreateTemp(dir, pattern string) (IFile, error)
	NewFile(fd uintptr, name string) IFile
	Open(name string) (IFile, error)
	OpenFile(name string, flag int, perm os.FileMode) (IFile, error)
	SetReadingFiles(files []IFile)
	SetWritingFiles(files []IFile)
	GetReadingFiles() []IFile
	GetWritingFiles() []IFile
	GetFiles() []IFile
	HasFile(file IFile) bool
	RemoveFile(file IFile) bool

	// similar to IFile
	Close() error
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	ReadFrom(r io.Reader) (n int64, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	Write(b []byte) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	WriteString(s string) (n int, err error)
	WriteTo(w io.Writer) (n int64, err error)
}
