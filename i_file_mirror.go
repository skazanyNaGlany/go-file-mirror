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
	SetReadingFiles(files []IFile) error
	SetWritingFiles(files []IFile) error
	GetReadingFiles() []IFile
	GetWritingFiles() []IFile
	GetFiles() []IFile
	HasFile(file IFile) bool
	RemoveFile(file IFile) bool

	// similar to IFile
	Close(file IFile) error
	Read(file IFile, b []byte) (n int, err error)
	ReadAt(file IFile, b []byte, off int64) (n int, err error)
	ReadFrom(file IFile, r io.Reader) (n int64, err error)
	Seek(file IFile, offset int64, whence int) (ret int64, err error)
	Stat(file IFile) (os.FileInfo, error)
	Sync(file IFile) error
	Truncate(file IFile, size int64) error
	Write(file IFile, b []byte) (n int, err error)
	WriteAt(file IFile, b []byte, off int64) (n int, err error)
	WriteString(file IFile, s string) (n int, err error)
	WriteTo(file IFile, w io.Writer) (n int64, err error)
}
