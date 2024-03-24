package gofilemirror

import (
	"io"
	"os"
)

// implements IFile, IFileEx
type File struct {
	fileMirror     IFileMirror
	underlyingFile *os.File
}

func NewFile(fileMirror IFileMirror, underlyingFile *os.File) *File {
	fmf := File{}

	fmf.fileMirror = fileMirror
	fmf.underlyingFile = underlyingFile

	return &fmf
}

// IFile
func (f *File) Close() error {
	return f.fileMirror.Close()
}

func (f *File) Read(b []byte) (n int, err error) {
	return f.fileMirror.Read(b)
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (f *File) ReadFrom(r io.Reader) (n int64, err error) {
	panic("not implemented")
}

func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	return f.fileMirror.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.fileMirror.Stat()
}

func (f *File) Sync() error {
	panic("not implemented")
}

func (f *File) Truncate(size int64) error {
	panic("not implemented")
}

func (f *File) Write(b []byte) (n int, err error) {
	return f.fileMirror.Write(b)
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (f *File) WriteString(s string) (n int, err error) {
	return f.fileMirror.WriteString(s)
}

func (f *File) WriteTo(w io.Writer) (n int64, err error) {
	panic("not implemented")
}

func (f *File) SetFileMirror(fileMirror IFileMirror) {
	f.fileMirror = fileMirror
}

func (f *File) SetUnderlyingFile(underlyingFile *os.File) {
	f.underlyingFile = underlyingFile
}

func (f *File) GetFileMirror() IFileMirror {
	return f.fileMirror
}

func (f *File) GetUnderlyingFile() *os.File {
	return f.underlyingFile
}
