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
	osf, err := os.CreateTemp(dir, pattern)

	fmf := NewFile(fm, osf)

	fm.readingFiles = append(fm.readingFiles, fmf)
	fm.writingFiles = append(fm.writingFiles, fmf)

	return fmf, err
}

func (fm *FileMirror) NewFile(fd uintptr, name string) IFile {
	osf := os.NewFile(fd, name)

	fmf := NewFile(fm, osf)

	fm.readingFiles = append(fm.readingFiles, fmf)

	return fmf
}

func (fm *FileMirror) Open(name string) (IFile, error) {
	osf, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	fmf := NewFile(fm, osf)

	fm.readingFiles = append(fm.readingFiles, fmf)

	return fmf, nil
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
	files := make([]IFile, 0)

	for _, f := range append(fm.readingFiles, fm.writingFiles...) {
		if !slices.Contains(files, f) {
			files = append(files, f)
		}
	}

	return files
}

func (fm *FileMirror) HasFile(file IFile) bool {
	for _, f := range fm.GetFiles() {
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

func (fm *FileMirror) Close() error {
	files := fm.GetFiles()

	if len(files) == 0 {
		return ErrNoFiles
	}

	for _, f := range files {
		if err := f.GetUnderlyingFile().Close(); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileMirror) Read(b []byte) (n int, err error) {
	if len(fm.readingFiles) == 0 {
		return 0, ErrNoFilesToRead
	}

	return fm.readingFiles[0].GetUnderlyingFile().Read(b)
}

func (fm *FileMirror) ReadAt(b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) ReadFrom(r io.Reader) (n int64, err error) {
	panic("not implemented")
}

func (fm *FileMirror) Seek(offset int64, whence int) (ret int64, err error) {
	files := fm.GetFiles()

	if len(files) == 0 {
		return 0, ErrNoFiles
	}

	for _, f := range files {
		ret, err = f.GetUnderlyingFile().Seek(offset, whence)

		if err != nil {
			return ret, err
		}
	}

	return ret, err
}

func (fm *FileMirror) Stat() (os.FileInfo, error) {
	files := fm.GetFiles()

	if len(files) == 0 {
		return nil, ErrNoFiles
	}

	return files[0].GetUnderlyingFile().Stat()
}

func (fm *FileMirror) Sync() error {
	panic("not implemented")
}

func (fm *FileMirror) Truncate(size int64) error {
	panic("not implemented")
}

func (fm *FileMirror) Write(b []byte) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) WriteAt(b []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (fm *FileMirror) WriteString(s string) (n int, err error) {
	if len(fm.writingFiles) == 0 {
		return 0, ErrNoFilesToWrite
	}

	for _, f := range fm.writingFiles {
		n, err = f.GetUnderlyingFile().WriteString(s)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (fm *FileMirror) WriteTo(w io.Writer) (n int64, err error) {
	panic("not implemented")
}

func NewFileMirror() IFileMirror {
	fm := FileMirror{}

	return &fm
}
