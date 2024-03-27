package gofilemirror

import (
	"os"
	"slices"
)

// implements IFileMirror
// driver for IFile
type FileMirror struct {
	readFiles  []IFile
	writeFiles []IFile
}

func (fm *FileMirror) AddReadingFile(file IFile) bool {
	if slices.Contains(fm.readFiles, file) {
		return false
	}

	fm.readFiles = append(fm.readFiles, file)
	file.SetFileMirror(fm)

	return true
}

func (fm *FileMirror) RemoveReadingFile(file IFile) bool {
	i := slices.Index(fm.readFiles, file)

	if i == -1 {
		return false
	}

	fm.readFiles = slices.Delete(fm.readFiles, i, i+1)
	file.SetFileMirror(nil)

	return true
}

func (fm *FileMirror) GetReadingFiles() []IFile {
	return fm.readFiles
}

func (fm *FileMirror) AddWritingFile(file IFile) bool {
	if slices.Contains(fm.writeFiles, file) {
		return false
	}

	fm.writeFiles = append(fm.writeFiles, file)
	file.SetFileMirror(fm)

	return true
}

func (fm *FileMirror) RemoveWritingFile(file IFile) bool {
	i := slices.Index(fm.writeFiles, file)

	if i == -1 {
		return false
	}

	fm.writeFiles = slices.Delete(fm.writeFiles, i, i+1)
	file.SetFileMirror(nil)

	return true

}

func (fm *FileMirror) GetWritingFiles() []IFile {
	return fm.writeFiles
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
	if len(fm.readFiles) == 0 {
		return 0, ErrNoFilesToRead
	}

	return fm.readFiles[0].GetUnderlyingFile().Read(b)
}

func (fm *FileMirror) ReadAt(b []byte, off int64) (n int, err error) {
	if len(fm.readFiles) == 0 {
		return 0, ErrNoFilesToRead
	}

	return fm.readFiles[0].GetUnderlyingFile().ReadAt(b, off)
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
	files := fm.GetFiles()

	if len(files) == 0 {
		return ErrNoFiles
	}

	for _, f := range files {
		if err := f.GetUnderlyingFile().Sync(); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileMirror) Truncate(size int64) error {
	if len(fm.writeFiles) == 0 {
		return ErrNoFilesToWrite
	}

	for _, f := range fm.writeFiles {
		if err := f.GetUnderlyingFile().Truncate(size); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileMirror) Write(b []byte) (n int, err error) {
	if len(fm.writeFiles) == 0 {
		return 0, ErrNoFilesToWrite
	}

	for _, f := range fm.writeFiles {
		n, err = f.GetUnderlyingFile().Write(b)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (fm *FileMirror) WriteAt(b []byte, off int64) (n int, err error) {
	if len(fm.writeFiles) == 0 {
		return 0, ErrNoFilesToWrite
	}

	for _, f := range fm.writeFiles {
		n, err = f.GetUnderlyingFile().WriteAt(b, off)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (fm *FileMirror) WriteString(s string) (n int, err error) {
	if len(fm.writeFiles) == 0 {
		return 0, ErrNoFilesToWrite
	}

	for _, f := range fm.writeFiles {
		n, err = f.GetUnderlyingFile().WriteString(s)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// TODO make private
func (fm *FileMirror) GetFiles() []IFile {
	files := make([]IFile, 0)

	for _, f := range append(fm.readFiles, fm.writeFiles...) {
		if !slices.Contains(files, f) {
			files = append(files, f)
		}
	}

	return files
}

func NewFileMirror() IFileMirror {
	fm := FileMirror{}

	return &fm
}
