package gofilemirror

import (
	"os"
	"slices"
)

// implements IFileMirror
// driver for IFile
type FileMirror struct {
	readFiles              []IFile
	writeFiles             []IFile
	asyncFiles             []IFile
	running                bool
	operations             chan *AsyncOperation
	asyncOperationCallback AsyncOperationCallback
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

func (fm *FileMirror) innerClose() error {
	fm.running = false

	if fm.operations != nil {
		close(fm.operations)
		fm.operations = nil
	}

	return nil
}

func (fm *FileMirror) AddAsyncFile(file IFile) bool {
	if slices.Contains(fm.asyncFiles, file) {
		return false
	}

	fm.asyncFiles = append(fm.asyncFiles, file)
	file.SetFileMirror(fm)

	return true
}

func (fm *FileMirror) RemoveAsyncFile(file IFile) bool {
	i := slices.Index(fm.asyncFiles, file)

	if i == -1 {
		return false
	}

	fm.asyncFiles = slices.Delete(fm.asyncFiles, i, i+1)
	file.SetFileMirror(nil)

	return true
}

func (fm *FileMirror) GetAsyncFiles() []IFile {
	return fm.readFiles
}

func (fm *FileMirror) Close() error {
	fm.innerClose()
	return nil
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

	fm.innerClose()
}

func (fm *FileMirror) execute(operation *AsyncOperation) {
	// TODO
	switch operation._type {
	case AOT_READ:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		n, err := operation.file.GetUnderlyingFile().Read(operation.buffer)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_READ_AT:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		n, err := operation.file.GetUnderlyingFile().
			ReadAt(operation.buffer, operation.offset)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_WRITE:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		n, err := operation.file.GetUnderlyingFile().Write(operation.buffer)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_WRITE_AT:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		n, err := operation.file.GetUnderlyingFile().
			WriteAt(operation.buffer, operation.offset)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_WRITE_STRING:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		n, err := operation.file.GetUnderlyingFile().
			WriteString(operation.stringBuffer)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_TRUNCATE:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		err := operation.file.GetUnderlyingFile().Truncate(operation.size)

		operation.err = err
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	case AOT_SEEK:
		if mutex := operation.file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}

		ret, err := operation.file.GetUnderlyingFile().
			Seek(operation.offset, operation.whence)

		operation.err = err
		operation.resultInt = ret
		operation.done = true

		if fm.asyncOperationCallback != nil {
			fm.asyncOperationCallback(operation)
		}
	default:
		panic("not implemented")
	}
}

func (fm *FileMirror) close() error {
	files := fm.getFiles()

	if len(files) == 0 {
		return ErrNoFiles
	}

	for _, f := range files {
		if mutex := f.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		if err := f.GetUnderlyingFile().Close(); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileMirror) read(b []byte) (operations []*AsyncOperation, n int, err error) {
	if len(fm.readFiles) == 0 {
		return nil, 0, ErrNoFilesToRead
	}

	file := fm.readFiles[0]

	if slices.Contains(fm.asyncFiles, file) {
		asyncOp := AsyncOperation{}

		asyncOp._type = AOT_READ
		asyncOp.file = file
		asyncOp.buffer = make([]byte, len(b))

		operations = append(operations, &asyncOp)

		fm.operations <- &asyncOp

		return operations, 0, nil
	} else {
		if mutex := file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		n, err = file.GetUnderlyingFile().Read(b)

		return nil, n, err
	}
}

func (fm *FileMirror) readAt(
	b []byte,
	off int64,
) (operations []*AsyncOperation, n int, err error) {
	if len(fm.readFiles) == 0 {
		return nil, 0, ErrNoFilesToRead
	}

	file := fm.readFiles[0]

	if slices.Contains(fm.asyncFiles, file) {
		asyncOp := AsyncOperation{}

		asyncOp._type = AOT_READ_AT
		asyncOp.file = file
		asyncOp.buffer = make([]byte, len(b))
		asyncOp.offset = off

		operations = append(operations, &asyncOp)

		fm.operations <- &asyncOp

		return operations, 0, nil
	} else {
		if mutex := file.GetMutex(); mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		n, err = file.GetUnderlyingFile().ReadAt(b, off)

		return nil, n, err
	}
}

func (fm *FileMirror) seek(
	offset int64,
	whence int,
) (operations []*AsyncOperation, ret int64, err error) {
	files := fm.getFiles()

	if len(files) == 0 {
		return nil, 0, ErrNoFiles
	}

	for _, file := range files {
		if slices.Contains(fm.asyncFiles, file) {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_SEEK
			asyncOp.file = file
			asyncOp.offset = offset
			asyncOp.whence = whence

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := file.GetMutex(); mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			ret, err = file.GetUnderlyingFile().Seek(offset, whence)

			if err != nil {
				return operations, ret, err
			}
		}
	}

	return operations, ret, err
}

func (fm *FileMirror) stat() (os.FileInfo, error) {
	if len(fm.readFiles) == 0 {
		return nil, ErrNoFilesToRead
	}

	if mutex := fm.readFiles[0].GetMutex(); mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	return fm.readFiles[0].GetUnderlyingFile().Stat()
}

func (fm *FileMirror) sync() (operations []*AsyncOperation, err error) {
	files := fm.getFiles()

	if len(files) == 0 {
		return nil, ErrNoFiles
	}

	for _, file := range files {
		if slices.Contains(fm.asyncFiles, file) {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_SYNC
			asyncOp.file = file

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := file.GetMutex(); mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			if err := file.GetUnderlyingFile().Sync(); err != nil {
				return operations, err
			}
		}
	}

	return operations, nil
}

func (fm *FileMirror) truncate(size int64) (operations []*AsyncOperation, err error) {
	if len(fm.writeFiles) == 0 {
		return nil, ErrNoFilesToWrite
	}

	for _, file := range fm.writeFiles {
		if slices.Contains(fm.asyncFiles, file) {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_TRUNCATE
			asyncOp.file = file
			asyncOp.size = size

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := file.GetMutex(); mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			if err := file.GetUnderlyingFile().Truncate(size); err != nil {
				return operations, err
			}
		}
	}

	return operations, err
}

func (fm *FileMirror) write(b []byte) (operations []*AsyncOperation, n int, err error) {
	if len(fm.writeFiles) == 0 {
		return nil, 0, ErrNoFilesToWrite
	}

	for _, file := range fm.writeFiles {
		if slices.Contains(fm.asyncFiles, file) {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_WRITE
			asyncOp.file = file
			asyncOp.buffer = make([]byte, len(b))

			copy(asyncOp.buffer, b)

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := file.GetMutex(); mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			n, err = file.GetUnderlyingFile().Write(b)

			if err != nil {
				return operations, n, err
			}
		}
	}

	return operations, n, nil
}

func (fm *FileMirror) writeAt(
	b []byte,
	off int64,
) (operations []*AsyncOperation, n int, err error) {
	if len(fm.writeFiles) == 0 {
		return nil, 0, ErrNoFilesToWrite
	}

	for _, file := range fm.writeFiles {
		if slices.Contains(fm.asyncFiles, file) {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_WRITE_AT
			asyncOp.file = file
			asyncOp.buffer = make([]byte, len(b))
			asyncOp.offset = off

			copy(asyncOp.buffer, b)

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := file.GetMutex(); mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			n, err = file.GetUnderlyingFile().WriteAt(b, off)

			if err != nil {
				return operations, n, err
			}
		}

	}

	return operations, n, nil
}

func (fm *FileMirror) writeString(
	s string,
) (operations []*AsyncOperation, n int, err error) {
	if len(fm.writeFiles) == 0 {
		return nil, 0, ErrNoFilesToWrite
	}

	for _, file := range fm.writeFiles {
		if slices.Contains(fm.asyncFiles, file) {
			asyncOp := AsyncOperation{}

			asyncOp._type = AOT_WRITE_STRING
			asyncOp.file = file
			asyncOp.stringBuffer = s

			operations = append(operations, &asyncOp)

			fm.operations <- &asyncOp
		} else {
			if mutex := file.GetMutex(); mutex != nil {
				mutex.Lock()
				defer mutex.Unlock()
			}

			n, err = file.GetUnderlyingFile().WriteString(s)

			if err != nil {
				return operations, n, err
			}
		}
	}

	return operations, n, nil
}

func (fm *FileMirror) getFiles() []IFile {
	files := make([]IFile, 0)

	for _, f := range append(fm.readFiles, fm.writeFiles...) {
		if !slices.Contains(files, f) {
			files = append(files, f)
		}
	}

	return files
}

func NewFileMirror(queueSize int) IFileMirror {
	fm := FileMirror{}
	fm.operations = make(chan *AsyncOperation, queueSize)

	go fm.run()

	return &fm
}
