package gofilemirror

import "errors"

var ErrNoFileToRead = errors.New("no files to read")
var ErrNoFilesToWrite = errors.New("no files to write")
var ErrNoFiles = errors.New("no files")
