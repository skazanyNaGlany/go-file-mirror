package gofilemirror

import (
	"os"
	"slices"
	"sync"
)

// implements IFileMirror
// driver for IFile
type FileMirror struct {
	readingFile            *os.File
	writingFiles           []*os.File
	fileMutexes            map[*os.File]*sync.Mutex
	asyncFiles             map[*os.File]bool
	running                bool
	operations             chan *AsyncOperation
	asyncOperationCallback AsyncOperationCallback
}

func (fm *FileMirror) SetFileMutex(file *os.File, mutex *sync.Mutex) {
	if mutex != nil {
		fm.fileMutexes[file] = mutex
	} else {
		delete(fm.fileMutexes, file)
	}
}

func (fm *FileMirror) SetFileAsync(file *os.File, async bool) {
	if async {
		fm.asyncFiles[file] = true
	} else {
		delete(fm.asyncFiles, file)
	}
}

func (fm *FileMirror) GetFileMutex(file *os.File) *sync.Mutex {
	return fm.fileMutexes[file]
}

func (fm *FileMirror) IsFileAsync(file *os.File) bool {
	return fm.asyncFiles[file]
}

func (fm *FileMirror) SetReadingFile(file *os.File) {
	fm.readingFile = file
}

func (fm *FileMirror) GetReadingFile() *os.File {
	return fm.readingFile
}

func (fm *FileMirror) AddWritingFile(file *os.File) bool {
	if slices.Contains(fm.writingFiles, file) {
		return false
	}

	fm.writingFiles = append(fm.writingFiles, file)

	return true
}

func (fm *FileMirror) RemoveWritingFile(file *os.File) bool {
	i := slices.Index(fm.writingFiles, file)

	if i == -1 {
		return false
	}

	fm.writingFiles = slices.Delete(fm.writingFiles, i, i+1)

	return true
}

func (fm *FileMirror) GetWritingFiles() []*os.File {
	return fm.writingFiles
}

func (fm *FileMirror) SetAsyncOperationCallback(callback AsyncOperationCallback) {
	fm.asyncOperationCallback = callback
}

func (fm *FileMirror) GetAsyncOperationCallback() AsyncOperationCallback {
	return fm.asyncOperationCallback
}

func (fm *FileMirror) run() {
	fm.running = true

	for asyncOp := range fm.operations {
		if asyncOp._type != AOT_NONE {
			fm.execute(asyncOp)
		}

		if !fm.running {
			break
		}
	}
}

func (fm *FileMirror) execute(operation *AsyncOperation) {
	switch operation._type {
	case AOT_READ:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			if !fm.asyncOperationCallback(operation) {
				return
			}
		}

		n, err := operation.file.Read(operation.buffer)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_READ_AT:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		n, err := operation.file.ReadAt(operation.buffer, operation.offset)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_WRITE:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		n, err := operation.file.Write(operation.buffer)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_WRITE_AT:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		n, err := operation.file.WriteAt(operation.buffer, operation.offset)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_WRITE_STRING:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		n, err := operation.file.WriteString(operation.stringBuffer)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_TRUNCATE:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		err := operation.file.Truncate(operation.size)

		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_SEEK:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		ret, err := operation.file.Seek(operation.offset, operation.whence)

		operation.err = err
		operation.resultInt = ret
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_SYNC:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if !fm.asyncOperationCallback(operation) {
			return
		}

		err := operation.file.Sync()

		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	}
}

func (fm *FileMirror) Close() error {
	fm.running = false

	if fm.operations != nil {
		close(fm.operations)
		fm.operations = nil
	}

	files := fm.GetAllFiles()

	if len(files) == 0 {
		return ErrNoFiles
	}

	for _, f := range files {
		if mutex := fm.fileMutexes[f]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileMirror) RemoveAllFiles() error {
	files := fm.GetAllFiles()

	if len(files) == 0 {
		return ErrNoFiles
	}

	fm.readingFile = nil
	fm.writingFiles = make([]*os.File, 0)
	fm.asyncFiles = make(map[*os.File]bool)
	fm.fileMutexes = make(map[*os.File]*sync.Mutex)

	return nil
}

func (fm *FileMirror) Read(
	b []byte,
	asyncOpUserData any,
) (operations []*AsyncOperation, n int, err error) {
	if fm.readingFile == nil {
		return nil, 0, ErrNoFileToRead
	}

	file := fm.readingFile

	if fm.asyncFiles[file] {
		asyncOp := AsyncOperation{}

		asyncOp._type = AOT_READ
		asyncOp.file = file
		asyncOp.buffer = make([]byte, len(b))
		asyncOp.userData = asyncOpUserData

		operations = append(operations, &asyncOp)

		fm.operations <- &asyncOp

		return operations, 0, nil
	} else {
		if mutex := fm.fileMutexes[file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		n, err = file.Read(b)

		return nil, n, err
	}
}

func (fm *FileMirror) ReadAt(
	b []byte,
	off int64,
	asyncOpUserData any,
) (operations []*AsyncOperation, n int, err error) {
	if fm.readingFile == nil {
		return nil, 0, ErrNoFileToRead
	}

	file := fm.readingFile

	if fm.asyncFiles[file] {
		asyncOp := AsyncOperation{}

		asyncOp._type = AOT_READ_AT
		asyncOp.file = file
		asyncOp.buffer = make([]byte, len(b))
		asyncOp.offset = off
		asyncOp.userData = asyncOpUserData

		operations = append(operations, &asyncOp)

		fm.operations <- &asyncOp

		return operations, 0, nil
	} else {
		if mutex := fm.fileMutexes[file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		n, err = file.ReadAt(b, off)

		return nil, n, err
	}
}

func (fm *FileMirror) Seek(
	offset int64,
	whence int,
	asyncOpUserData any,
) (operations []*AsyncOperation, ret int64, err error) {
	files := fm.GetAllFiles()

	if len(files) == 0 {
		return nil, 0, ErrNoFiles
	}

	for _, file := range files {
		if fm.asyncFiles[file] {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_SEEK
			asyncOp.file = file
			asyncOp.offset = offset
			asyncOp.whence = whence
			asyncOp.userData = asyncOpUserData

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := fm.fileMutexes[file]; mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			ret, err = file.Seek(offset, whence)

			if err != nil {
				return operations, ret, err
			}
		}
	}

	return operations, ret, err
}

func (fm *FileMirror) Stat() (os.FileInfo, error) {
	if fm.readingFile == nil {
		return nil, ErrNoFileToRead
	}

	if mutex := fm.fileMutexes[fm.readingFile]; mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	return fm.readingFile.Stat()
}

func (fm *FileMirror) Sync(
	asyncOpUserData any,
) (operations []*AsyncOperation, err error) {
	files := fm.GetAllFiles()

	if len(files) == 0 {
		return nil, ErrNoFiles
	}

	for _, file := range files {
		if fm.asyncFiles[file] {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_SYNC
			asyncOp.file = file
			asyncOp.userData = asyncOpUserData

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := fm.fileMutexes[file]; mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			if err := file.Sync(); err != nil {
				return operations, err
			}
		}
	}

	return operations, nil
}

func (fm *FileMirror) Truncate(
	size int64,
	asyncOpUserData any,
) (operations []*AsyncOperation, err error) {
	if len(fm.writingFiles) == 0 {
		return nil, ErrNoFilesToWrite
	}

	for _, file := range fm.writingFiles {
		if fm.asyncFiles[file] {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_TRUNCATE
			asyncOp.file = file
			asyncOp.size = size
			asyncOp.userData = asyncOpUserData

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := fm.fileMutexes[file]; mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			if err := file.Truncate(size); err != nil {
				return operations, err
			}
		}
	}

	return operations, err
}

func (fm *FileMirror) Write(
	b []byte,
	asyncOpUserData any,
) (operations []*AsyncOperation, n int, err error) {
	if len(fm.writingFiles) == 0 {
		return nil, 0, ErrNoFilesToWrite
	}

	for _, file := range fm.writingFiles {
		if fm.asyncFiles[file] {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_WRITE
			asyncOp.file = file
			asyncOp.buffer = make([]byte, len(b))
			asyncOp.userData = asyncOpUserData

			copy(asyncOp.buffer, b)

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := fm.fileMutexes[file]; mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			n, err = file.Write(b)

			if err != nil {
				return operations, n, err
			}
		}
	}

	return operations, n, nil
}

func (fm *FileMirror) WriteAt(
	b []byte,
	off int64,
	asyncOpUserData any,
) (operations []*AsyncOperation, n int, err error) {
	if len(fm.writingFiles) == 0 {
		return nil, 0, ErrNoFilesToWrite
	}

	for _, file := range fm.writingFiles {
		if fm.asyncFiles[file] {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_WRITE_AT
			asyncOp.file = file
			asyncOp.buffer = make([]byte, len(b))
			asyncOp.offset = off
			asyncOp.userData = asyncOpUserData

			copy(asyncOp.buffer, b)

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := fm.fileMutexes[file]; mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			n, err = file.WriteAt(b, off)

			if err != nil {
				return operations, n, err
			}
		}

	}

	return operations, n, nil
}

func (fm *FileMirror) WriteString(
	s string,
	asyncOpUserData any,
) (operations []*AsyncOperation, n int, err error) {
	if len(fm.writingFiles) == 0 {
		return nil, 0, ErrNoFilesToWrite
	}

	for _, file := range fm.writingFiles {
		if fm.asyncFiles[file] {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_WRITE_STRING
			asyncOp.file = file
			asyncOp.stringBuffer = s
			asyncOp.userData = asyncOpUserData

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := fm.fileMutexes[file]; mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			n, err = file.WriteString(s)

			if err != nil {
				return operations, n, err
			}
		}
	}

	return operations, n, nil
}

func (fm *FileMirror) GetAllFiles() []*os.File {
	files := slices.Clone(fm.writingFiles)

	if fm.readingFile == nil {
		return files
	}

	if !slices.Contains(files, fm.readingFile) {
		files = append(files, fm.readingFile)
	}

	return files
}

func NewFileMirror(queueSize int) *FileMirror {
	fm := FileMirror{}
	fm.operations = make(chan *AsyncOperation, queueSize)
	fm.fileMutexes = make(map[*os.File]*sync.Mutex)
	fm.asyncFiles = make(map[*os.File]bool)

	go fm.run()

	return &fm
}
