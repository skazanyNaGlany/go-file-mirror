package gofilemirror

import (
	"os"
)

// implements IFile
type IFileMirror interface {
	AddReadingFile(file IFile) bool
	RemoveReadingFile(file IFile) bool
	GetReadingFiles() []IFile
	AddWritingFile(file IFile) bool
	RemoveWritingFile(file IFile) bool
	GetWritingFiles() []IFile

	// similar to IFile
	Close() error
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	Write(b []byte) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	WriteString(s string) (n int, err error)
}
