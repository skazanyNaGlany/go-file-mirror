package gofilemirror

import (
	"os"
	"sync"
)

type IFile interface {
	GetFileMirror() IFileMirror
	SetFileMirror(fileMirror IFileMirror)
	GetMutex() *sync.Mutex
	SetMutex(mutex *sync.Mutex)
	GetUnderlyingFile() *os.File
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
