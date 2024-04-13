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
	Close() error
	SetAsyncOperationCallback(callback AsyncOperationCallback)
	GetAsyncOperationCallback() AsyncOperationCallback

	// similar to IFile
	close() error
	read(b []byte) (n int, err error)
	readAt(b []byte, off int64) (n int, err error)
	seek(offset int64, whence int) (ret int64, err error)
	stat() (os.FileInfo, error)
	sync() error
	truncate(size int64) error
	write(b []byte) (n int, err error)
	writeAt(b []byte, off int64) (n int, err error)
	writeString(s string) (n int, err error)
	readAsync(b []byte) *AsyncOperation
	readAtAsync(b []byte, off int64) *AsyncOperation
	seekAsync(offset int64, whence int) *AsyncOperation
	statAsync() *AsyncOperation
	syncAsync() *AsyncOperation
	truncateAsync(size int64) *AsyncOperation
	writeAsync(b []byte) *AsyncOperation
	writeAtAsync(b []byte, off int64) *AsyncOperation
	writeStringAsync(s string) *AsyncOperation
}
