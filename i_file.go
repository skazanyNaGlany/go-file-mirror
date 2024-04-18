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
	Read(b []byte) (ops []*AsyncOperation, n int, err error)
	ReadAt(b []byte, off int64) (ops []*AsyncOperation, n int, err error)
	Seek(offset int64, whence int) (ops []*AsyncOperation, ret int64, err error)
	Stat() (os.FileInfo, error)
	Sync() (ops []*AsyncOperation, err error)
	Truncate(size int64) (ops []*AsyncOperation, err error)
	Write(b []byte) (ops []*AsyncOperation, n int, err error)
	WriteAt(b []byte, off int64) (ops []*AsyncOperation, n int, err error)
	WriteString(s string) (ops []*AsyncOperation, n int, err error)
}
