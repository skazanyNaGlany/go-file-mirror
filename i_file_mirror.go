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
	AddAsyncFile(file IFile) bool
	RemoveAsyncFile(file IFile) bool
	GetAsyncFiles() []IFile
	Close() error
	SetAsyncOperationCallback(callback AsyncOperationCallback)
	GetAsyncOperationCallback() AsyncOperationCallback

	// similar to IFile
	close() error
	read(b []byte) (ops []*AsyncOperation, n int, err error)
	readAt(b []byte, off int64) (ops []*AsyncOperation, n int, err error)
	seek(offset int64, whence int) (ops []*AsyncOperation, ret int64, err error)
	stat() (os.FileInfo, error)
	sync() (ops []*AsyncOperation, err error)
	truncate(size int64) (ops []*AsyncOperation, err error)
	write(b []byte) (ops []*AsyncOperation, n int, err error)
	writeAt(b []byte, off int64) (ops []*AsyncOperation, n int, err error)
	writeString(s string) (ops []*AsyncOperation, n int, err error)
}
