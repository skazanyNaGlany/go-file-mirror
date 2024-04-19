package gofilemirror

import (
	"os"
	"sync"
)

// implements IFile
type File struct {
	fileMirror     IFileMirror
	underlyingFile *os.File
	mutex          *sync.Mutex
}

// IFile
func (f *File) Close() error {
	return f.fileMirror.close()
}

func (f *File) Read(b []byte) (operations []*AsyncOperation, n int, err error) {
	return f.fileMirror.read(b)
}

func (f *File) ReadAt(
	b []byte,
	off int64,
) (operations []*AsyncOperation, n int, err error) {
	return f.fileMirror.readAt(b, off)
}

func (f *File) Seek(
	offset int64,
	whence int,
) (operations []*AsyncOperation, ret int64, err error) {
	return f.fileMirror.seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.fileMirror.stat()
}

func (f *File) Sync() (operations []*AsyncOperation, err error) {
	return f.fileMirror.sync()
}

func (f *File) Truncate(size int64) (operations []*AsyncOperation, err error) {
	return f.fileMirror.truncate(size)
}

func (f *File) Write(b []byte) (operations []*AsyncOperation, n int, err error) {
	return f.fileMirror.write(b)
}

func (f *File) WriteAt(
	b []byte,
	off int64,
) (operations []*AsyncOperation, n int, err error) {
	return f.fileMirror.writeAt(b, off)
}

func (f *File) WriteString(s string) (operations []*AsyncOperation, n int, err error) {
	return f.fileMirror.writeString(s)
}

func (f *File) GetFileMirror() IFileMirror {
	return f.fileMirror
}

func (f *File) SetFileMirror(fileMirror IFileMirror) {
	f.fileMirror = fileMirror
}

func (f *File) GetUnderlyingFile() *os.File {
	return f.underlyingFile
}

func (f *File) GetMutex() *sync.Mutex {
	return f.mutex
}

func (f *File) SetMutex(mutex *sync.Mutex) {
	f.mutex = mutex
}

// Globals
func Create(name string) (IFile, error) {
	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f, mutex: &sync.Mutex{}}, nil
}

func CreateTemp(dir, pattern string) (IFile, error) {
	f, err := os.CreateTemp(dir, pattern)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f, mutex: &sync.Mutex{}}, nil
}

func NewFile(fd uintptr, name string) IFile {
	f := os.NewFile(fd, name)

	if f == nil {
		return nil
	}

	return &File{underlyingFile: f, mutex: &sync.Mutex{}}
}

func Open(name string) (IFile, error) {
	f, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f, mutex: &sync.Mutex{}}, nil
}

func OpenFile(name string, flag int, perm os.FileMode) (IFile, error) {
	f, err := os.OpenFile(name, flag, perm)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f, mutex: &sync.Mutex{}}, nil
}
