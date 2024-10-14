# GO File Mirror

Operate on many files as if it were one file.

## Overview

GO File Mirror allows to use many files as if it were one file, it supports reading and writing files, also asynchronously.

## Why GO File Mirror

Consider you have two or more files, one of them is stored on a fast medium, other is on a slow medium, and you want to perform I/O operations on these files as if it were one file, but want to read only from file on the fast medium.

## Installation

To install GO File Mirror use `go get`:

```
$ go get github.com/skazanyNaGlany/go-file-mirror
```

## Usage
```go
package main

import (
    "crypto/rand"
    "log"
    "os"
    "time"

    gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
)

func main() {
    // Create FileMirror instance
    fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
    defer fm.Close(true)

    fm.SetOperationCallback(func(operation *gofilemirror.Operation) {
        // The callback will be fired on any I/O operation
        // for each of the file handle
    })

    fm.SetIdleCallback(func(fileMirror *gofilemirror.FileMirror) {
        // The callback will be fired on idle FileMirror instance
    })

    // First temporary file
    f, err := os.CreateTemp("/tmp", "testing_file_mirror")
    if err != nil {
        panic(err)
    }

    // Second temporary file
    f2, err := os.CreateTemp("/tmp", "testing_file_mirror2")
    if err != nil {
        panic(err)
    }

    // 10 bytes buffer for I/O operations
    buffer := make([]byte, 10)

    // Fill the buffer with some random data
    rand.Read(buffer)

    // Reading only from "f"
    // Writing to "f" and "f2"
    // All operations on both files will be async
    fm.SetFileAsync(f, true)
    fm.SetFileAsync(f2, true)

    // Write 10 bytes at 0 offset to both files
    // The operation will be async, will not block
    // and return immediately
    operations := fm.WriteAt(buffer, 0)
    operations.WaitForDone(10 * time.Second)

    if len(operations.GetPendingOperations()) > 0 {
        log.Fatal("still have some pending operations")
    }

    // Fatal when one of the operations failed
    for _, operation := range *operations {
        if err := operation.GetLastResultError(); err != nil {
            log.Fatal(err)
        }
    }

    // Read written data
    // The operation will be async, will not block
    // and return immediately
    operations = fm.ReadAt(buffer, 0)
    operations.WaitForDone(10 * time.Second)

    if len(operations.GetPendingOperations()) > 0 {
        log.Fatal("still have some pending operations")
    }

    // Fatal when one of the operations failed
    for _, operation := range *operations {
        if err := operation.GetLastResultError(); err != nil {
            log.Fatal(err)
        }
    }

    firstOperation := operations.GetFirstAsyncOperation()

    log.Println(firstOperation.GetBuffer())
}
```

## License

© Paweł Kacperski, 2024 ~ time.Now

Released under the [MIT License](https://github.com/go-gorm/gorm/blob/master/LICENSE)
