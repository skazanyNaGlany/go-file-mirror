package gofilemirror

import (
	"os"
	"slices"
	"sync"
	"time"
)

type FileMirror struct {
	readingFiles      []*os.File
	writingFiles      []*os.File
	fileMutexes       map[*os.File]*sync.Mutex
	asyncFiles        map[*os.File]bool
	fileUserData      map[*os.File]any
	running           bool
	asyncOperations   chan *Operation
	operationCallback OperationCallback
	fixedBuffer       bool
}

func (fm *FileMirror) SetFixedBuffer(fixedBuffer bool) {
	fm.fixedBuffer = fixedBuffer
}

func (fm *FileMirror) IsFixedBuffer() bool {
	return fm.fixedBuffer
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

func (fm *FileMirror) SetFileUserData(file *os.File, userData any) {
	if userData != nil {
		fm.fileUserData[file] = userData
	} else {
		delete(fm.fileUserData, file)
	}
}

func (fm *FileMirror) GetFileMutex(file *os.File) *sync.Mutex {
	return fm.fileMutexes[file]
}

func (fm *FileMirror) IsFileAsync(file *os.File) bool {
	return fm.asyncFiles[file]
}

func (fm *FileMirror) GetFileUserData(file *os.File) any {
	return fm.fileUserData[file]
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

func (fm *FileMirror) AddReadingFile(file *os.File) bool {
	if slices.Contains(fm.readingFiles, file) {
		return false
	}

	fm.readingFiles = append(fm.readingFiles, file)

	return true
}

func (fm *FileMirror) RemoveReadingFile(file *os.File) bool {
	i := slices.Index(fm.readingFiles, file)

	if i == -1 {
		return false
	}

	fm.readingFiles = slices.Delete(fm.readingFiles, i, i+1)

	return true
}

func (fm *FileMirror) GetReadingFiles() []*os.File {
	return fm.readingFiles
}

func (fm *FileMirror) SetOperationCallback(callback OperationCallback) {
	fm.operationCallback = callback
}

func (fm *FileMirror) GetOperationCallback() OperationCallback {
	return fm.operationCallback
}

func (fm *FileMirror) run() {
	fm.running = true

	for asyncOperation := range fm.asyncOperations {
		if asyncOperation._type != OT_NONE {
			fm.execute(asyncOperation)
		}

		if !fm.running {
			break
		}
	}
}

func (fm *FileMirror) execute(operation *Operation) {
	switch operation._type {
	case OT_READ_AT:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.operationCallback != nil {
			fm.operationCallback(operation)
		}

		n, err := operation.file.ReadAt(operation.buffer, operation.offset)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.operationCallback != nil {
			fm.operationCallback(operation)
		}
	case OT_WRITE_AT:
		if mutex := fm.fileMutexes[operation.file]; mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}

		operation.started = true

		if fm.operationCallback != nil {
			fm.operationCallback(operation)
		}

		n, err := operation.file.WriteAt(operation.buffer, operation.offset)

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		if fm.operationCallback != nil {
			fm.operationCallback(operation)
		}
	}
}

func (fm *FileMirror) Close() error {
	fm.running = false

	if fm.asyncOperations != nil {
		close(fm.asyncOperations)
		fm.asyncOperations = nil
	}

	files := fm.GetAllFiles()

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
	fm.readingFiles = make([]*os.File, 0)
	fm.writingFiles = make([]*os.File, 0)
	fm.asyncFiles = make(map[*os.File]bool)
	fm.fileMutexes = make(map[*os.File]*sync.Mutex)

	return nil
}

func (fm *FileMirror) ReadAt(
	b []byte,
	off int64,
	operationUserData any,
	useFiles ...*os.File,
) (operationList *OperationList) {
	operationList = &OperationList{}

	if len(useFiles) == 0 {
		useFiles = fm.readingFiles
	}

	for _, file := range useFiles {
		operation := Operation{}

		operation._type = OT_READ_AT
		operation.file = file
		operation.offset = off
		operation.operationUserData = operationUserData
		operation.fileUserData = fm.fileUserData[file]

		if !fm.fixedBuffer {
			operation.buffer = make([]byte, len(b))
			copy(operation.buffer, b)
		} else {
			operation.buffer = b
		}

		*operationList = append(*operationList, &operation)

		if fm.asyncFiles[file] {
			operation.async = true

			fm.asyncOperations <- &operation
		} else {
			fm.execute(&operation)
		}
	}

	return operationList
}

func (fm *FileMirror) WriteAt(
	b []byte,
	off int64,
	operationUserData any,
	useFiles ...*os.File,
) (operationList *OperationList) {
	operationList = &OperationList{}

	if len(useFiles) == 0 {
		useFiles = fm.writingFiles
	}

	for _, file := range useFiles {
		operation := Operation{}

		operation._type = OT_WRITE_AT
		operation.file = file
		operation.offset = off
		operation.operationUserData = operationUserData
		operation.fileUserData = fm.fileUserData[file]

		if !fm.fixedBuffer {
			operation.buffer = make([]byte, len(b))
			copy(operation.buffer, b)
		} else {
			operation.buffer = b
		}

		*operationList = append(*operationList, &operation)

		if fm.asyncFiles[file] {
			operation.async = true

			fm.asyncOperations <- &operation
		} else {
			fm.execute(&operation)
		}
	}

	return operationList
}

func (fm *FileMirror) GetAllFiles() []*os.File {
	files := make([]*os.File, 0)

	for _, ifile := range fm.readingFiles {
		if !slices.Contains(files, ifile) {
			files = append(files, ifile)
		}
	}

	for _, ifile := range fm.writingFiles {
		if !slices.Contains(files, ifile) {
			files = append(files, ifile)
		}
	}

	return files
}

func (fm *FileMirror) WaitForNoAsyncOperations(duration time.Duration) {
	currentDuration := time.Duration(0)

	for {
		if len(fm.asyncOperations) == 0 {
			break
		}

		sleepDuration := 10 * time.Millisecond

		time.Sleep(sleepDuration)
		currentDuration += sleepDuration

		if currentDuration >= duration {
			break
		}
	}
}

func NewFileMirror(queueSize int) *FileMirror {
	fm := FileMirror{}
	fm.asyncOperations = make(chan *Operation, queueSize)
	fm.fileMutexes = make(map[*os.File]*sync.Mutex)
	fm.asyncFiles = make(map[*os.File]bool)
	fm.fileUserData = make(map[*os.File]any)

	go fm.run()

	return &fm
}
