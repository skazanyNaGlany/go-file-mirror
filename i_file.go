package gofilemirror

// IFile is an interface that represents a file.
// It is used to abstract the file operations.
type IFile interface {
	ReadAt(b []byte, off int64) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	Close() error
}
