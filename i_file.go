package gofilemirror

import (
	"io"
	"os"
)

type IFile interface {
	GetFileMirror() IFileMirror
	SetFileMirror(fileMirror IFileMirror)
	GetUnderlyingFile() *os.File
	Close() error
	Read(b []byte) (n int, err error)
	ReadFrom(r io.Reader) (n int64, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	WriteTo(w io.Writer) (n int64, err error)
}
