package gofilemirror

import (
	"io"
	"os"
)

// implements IFile
type File struct {
	fileMirror     IFileMirror
	underlyingFile *os.File
}

// IFile
func (f *File) Close() error {
	return f.fileMirror.Close()
}

func (f *File) Read(b []byte) (n int, err error) {
	return f.fileMirror.Read(b)
}

func (f *File) ReadFrom(r io.Reader) (n int64, err error) {
	return f.fileMirror.ReadFrom(r)
}

func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	return f.fileMirror.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.fileMirror.Stat()
}

func (f *File) Sync() error {
	return f.fileMirror.Sync()
}

func (f *File) Truncate(size int64) error {
	return f.fileMirror.Truncate(size)
}

func (f *File) Write(b []byte) (n int, err error) {
	return f.fileMirror.Write(b)
}

func (f *File) WriteString(s string) (n int, err error) {
	return f.fileMirror.WriteString(s)
}

func (f *File) WriteTo(w io.Writer) (n int64, err error) {
	return f.fileMirror.WriteTo(w)
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

// Globals
func Create(name string) (IFile, error) {
	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f}, nil
}

func CreateTemp(dir, pattern string) (IFile, error) {
	f, err := os.CreateTemp(dir, pattern)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f}, nil
}

func NewFile(fd uintptr, name string) IFile {
	f := os.NewFile(fd, name)

	if f == nil {
		return nil
	}

	return &File{underlyingFile: f}
}

func Open(name string) (IFile, error) {
	f, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f}, nil
}

func OpenFile(name string, flag int, perm os.FileMode) (IFile, error) {
	f, err := os.OpenFile(name, flag, perm)

	if err != nil {
		return nil, err
	}

	return &File{underlyingFile: f}, nil
}
