package gofilemirror

import (
	"io"
	"os"
	"slices"
)

// implements IFileMirror
// driver for IFile
type FileMirror struct {
	readingFiles []IFile
	writingFiles []IFile
}

func (fm *FileMirror) CreateTemp(dir, pattern string) (IFile, error) {
	f, err := os.CreateTemp(dir, pattern)

	fmf := NewFile(fm, f)

	fm.readingFiles = append(fm.readingFiles, fmf)
	fm.writingFiles = append(fm.writingFiles, fmf)

	return fmf, err
}

func (fm *FileMirror) NewFile(fd uintptr, name string) IFile {
	panic("not implemented")
}

func (fm *FileMirror) Open(name string) (IFile, error) {
	panic("not implemented")
}

func (fm *FileMirror) OpenFile(name string, flag int, perm os.FileMode) (IFile, error) {
	panic("not implemented")
}

func (fm *FileMirror) SetReadingFiles(files []IFile) error {
	panic("not implemented")
}

func (fm *FileMirror) SetWritingFiles(files []IFile) error {
	panic("not implemented")
}

func (fm *FileMirror) GetReadingFiles() []IFile {
	panic("not implemented")
}

func (fm *FileMirror) GetWritingFiles() []IFile {
	panic("not implemented")
}

func (fm *FileMirror) GetFiles() []IFile {
	panic("not implemented")
}

func (fm *FileMirror) HasFile(file IFile) bool {
	for _, f := range fm.readingFiles {
		if f == file {
			return true
		}
	}

	for _, f := range fm.writingFiles {
		if f == file {
			return true
		}
	}

	return false
}

func (fm *FileMirror) RemoveFile(file IFile) bool {
	for i, f := range fm.readingFiles {
		if f == file {
			fm.readingFiles = slices.Delete(fm.readingFiles, i, i+1)

			return true
		}
	}

	for i, f := range fm.writingFiles {
		if f == file {
			fm.writingFiles = slices.Delete(fm.writingFiles, i, i+1)

			return true
		}
	}

	return false
}

func (fm *FileMirror) Close(file IFile) error {
	if !fm.HasFile(file) {
		return ErrDoNotBelong
	}

	iFileEx := file.(IFileEx)

	err := iFileEx.GetUnderlyingFile().Close()
	iFileEx.SetUnderlyingFile(nil)

	fm.RemoveFile(file)

	return err
}

func (fm *FileMirror) Read(file IFile, b []byte) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) ReadAt(file IFile, b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) ReadFrom(file IFile, r io.Reader) (n int64, err error) {
	panic("not implemented")
}

func (fm *FileMirror) Seek(file IFile, offset int64, whence int) (ret int64, err error) {
	panic("not implemented")
}

func (fm *FileMirror) Stat(file IFile) (os.FileInfo, error) {
	panic("not implemented")
}

func (fm *FileMirror) Sync(file IFile) error {
	panic("not implemented")
}

func (fm *FileMirror) Truncate(file IFile, size int64) error {
	panic("not implemented")
}

func (fm *FileMirror) Write(file IFile, b []byte) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) WriteAt(file IFile, b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) WriteString(file IFile, s string) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) WriteTo(file IFile, w io.Writer) (n int64, err error) {
	panic("not implemented")
}

func NewFileMirror() IFileMirror {
	fm := FileMirror{}

	return &fm
}
