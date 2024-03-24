package gofilemirror

import "os"

type IFileEx interface {
	IFile

	SetFileMirror(fileMirror IFileMirror)
	SetUnderlyingFile(underlyingFile *os.File)
	GetFileMirror() IFileMirror
	GetUnderlyingFile() *os.File
}
