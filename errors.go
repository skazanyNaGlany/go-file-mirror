package gofilemirror

import "errors"

var ErrDoNotBelong = errors.New("file does not belong to FileMirror instance")
