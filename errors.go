package gofilemirror

import "errors"

var ErrDoNotBelong = errors.New("file does not belong to FileMirror instance")
var ErrNoFilesToRead = errors.New("no files to read")
var ErrNoFilesToWrite = errors.New("no files to write")
var ErrNoFiles = errors.New("no files")
