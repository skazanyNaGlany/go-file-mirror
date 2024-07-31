package gofilemirror

import (
	"slices"
	"sync"
	"time"
)

type FileMirror struct {
	readingFiles          []IFile
	writingFiles          []IFile
	allFiles              []IFile
	fileMutexes           map[IFile]*sync.Mutex
	asyncFiles            map[IFile]bool
	fileUserData          map[IFile]any
	running               bool
	asyncOperations       chan *Operation
	operationCallback     OperationCallback
	fixedBuffer           bool
	fileCachedMemoryBytes map[IFile][]bool
	idleCallback          IdleCallback
	idleSleepDuration     time.Duration
}

func (fm *FileMirror) GetFileCachedMemoryBytes(file IFile) []bool {
	return fm.fileCachedMemoryBytes[file]
}

func (fm *FileMirror) SetFileCachedMemoryBytes(file IFile, cachedMemoryBytes []bool) {
	if cachedMemoryBytes != nil {
		fm.fileCachedMemoryBytes[file] = cachedMemoryBytes
	} else {
		delete(fm.fileCachedMemoryBytes, file)
	}
}

func (fm *FileMirror) IsFileFullyCached(file IFile) bool {
	if fm.fileCachedMemoryBytes[file] == nil {
		return false
	}

	cachedBytes := 0

	for _, b := range fm.fileCachedMemoryBytes[file] {
		if b {
			cachedBytes++
		}
	}

	return cachedBytes == len(fm.fileCachedMemoryBytes[file])
}

func (fm *FileMirror) GetFileCachedPercent(file IFile) int {
	if fm.fileCachedMemoryBytes[file] == nil {
		return 0
	}

	maxLen := len(fm.fileCachedMemoryBytes[file])

	if maxLen == 0 {
		return 0
	}

	cachedBytes := 0

	for _, b := range fm.fileCachedMemoryBytes[file] {
		if b {
			cachedBytes++
		}
	}

	if cachedBytes == 0 {
		return 0
	}

	return cachedBytes * 100 / maxLen
}

func (fm *FileMirror) SetFixedBuffer(fixedBuffer bool) {
	fm.fixedBuffer = fixedBuffer
}

func (fm *FileMirror) IsFixedBuffer() bool {
	return fm.fixedBuffer
}

func (fm *FileMirror) SetFileMutex(file IFile, mutex *sync.Mutex) {
	if mutex != nil {
		fm.fileMutexes[file] = mutex
	} else {
		delete(fm.fileMutexes, file)
	}
}

func (fm *FileMirror) SetFileAsync(file IFile, async bool) {
	if async {
		fm.asyncFiles[file] = true
	} else {
		delete(fm.asyncFiles, file)
	}
}

func (fm *FileMirror) SetFileUserData(file IFile, userData any) {
	if userData != nil {
		fm.fileUserData[file] = userData
	} else {
		delete(fm.fileUserData, file)
	}
}

func (fm *FileMirror) GetFileMutex(file IFile) *sync.Mutex {
	return fm.fileMutexes[file]
}

func (fm *FileMirror) IsFileAsync(file IFile) bool {
	return fm.asyncFiles[file]
}

func (fm *FileMirror) GetFileUserData(file IFile) any {
	return fm.fileUserData[file]
}

func (fm *FileMirror) AddWritingFile(file IFile) bool {
	if slices.Contains(fm.writingFiles, file) {
		return false
	}

	fm.writingFiles = append(fm.writingFiles, file)

	if !slices.Contains(fm.allFiles, file) {
		fm.allFiles = append(fm.allFiles, file)
	}

	return true
}

func (fm *FileMirror) RemoveWritingFile(file IFile) bool {
	// writingFiles
	i := slices.Index(fm.writingFiles, file)

	if i == -1 {
		return false
	}

	fm.writingFiles = slices.Delete(fm.writingFiles, i, i+1)

	// allFiles
	i = slices.Index(fm.allFiles, file)

	if i != -1 {
		fm.allFiles = slices.Delete(fm.allFiles, i, i+1)
	}

	return true
}

func (fm *FileMirror) GetWritingFiles() []IFile {
	return fm.writingFiles
}

func (fm *FileMirror) AddReadingFile(file IFile) bool {
	if slices.Contains(fm.readingFiles, file) {
		return false
	}

	fm.readingFiles = append(fm.readingFiles, file)

	if !slices.Contains(fm.allFiles, file) {
		fm.allFiles = append(fm.allFiles, file)
	}

	return true
}

func (fm *FileMirror) RemoveReadingFile(file IFile) bool {
	// readingFiles
	i := slices.Index(fm.readingFiles, file)

	if i == -1 {
		return false
	}

	fm.readingFiles = slices.Delete(fm.readingFiles, i, i+1)

	// allFiles
	i = slices.Index(fm.allFiles, file)

	if i != -1 {
		fm.allFiles = slices.Delete(fm.allFiles, i, i+1)
	}

	return true
}

func (fm *FileMirror) GetReadingFiles() []IFile {
	return fm.readingFiles
}

func (fm *FileMirror) SetOperationCallback(callback OperationCallback) {
	fm.operationCallback = callback
}

func (fm *FileMirror) GetOperationCallback() OperationCallback {
	return fm.operationCallback
}

func (fm *FileMirror) SetIdleCallback(callback IdleCallback) {
	fm.idleCallback = callback
}

func (fm *FileMirror) GetIdleCallback() IdleCallback {
	return fm.idleCallback
}

func (fm *FileMirror) run() {
	fm.running = true

	for fm.running {
		select {
		case asyncOperation := <-fm.asyncOperations:
			if asyncOperation._type != OT_NONE {
				fm.execute(asyncOperation)
			}
		default:
			if fm.idleCallback != nil {
				fm.idleCallback(fm)
			}
		}

		if fm.idleSleepDuration > 0 {
			if len(fm.asyncOperations) == 0 {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (fm *FileMirror) fillCachedMemoryBytes(
	file IFile,
	startOffset int64,
	endOffset int64,
	b bool,
) {
	if fm.fileCachedMemoryBytes[file] == nil {
		return
	}

	maxLen := int64(len(fm.fileCachedMemoryBytes[file]))

	for i := startOffset; i < endOffset; i++ {
		if i < 0 || i > maxLen {
			return
		}

		fm.fileCachedMemoryBytes[file][i] = b
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

		startTime := time.Now().UnixMilli()

		n, err := operation.file.ReadAt(operation.buffer, operation.offset)

		operation.timeMilisecond = time.Now().UnixMilli() - startTime

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		fm.fillCachedMemoryBytes(
			operation.file,
			operation.offset,
			operation.offset+operation.resultInt,
			true,
		)

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

		startTime := time.Now().UnixMilli()

		n, err := operation.file.WriteAt(operation.buffer, operation.offset)

		operation.timeMilisecond = time.Now().UnixMilli() - startTime

		operation.resultInt = int64(n)
		operation.err = err
		operation.done = true

		fm.fillCachedMemoryBytes(
			operation.file,
			operation.offset,
			operation.offset+operation.resultInt,
			false,
		)

		if fm.operationCallback != nil {
			fm.operationCallback(operation)
		}
	}
}

func (fm *FileMirror) Close(closeOSHandles bool) error {
	fm.running = false

	if fm.asyncOperations != nil {
		close(fm.asyncOperations)
		fm.asyncOperations = nil
	}

	if !closeOSHandles {
		return nil
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
	fm.readingFiles = make([]IFile, 0)
	fm.writingFiles = make([]IFile, 0)
	fm.asyncFiles = make(map[IFile]bool)
	fm.fileMutexes = make(map[IFile]*sync.Mutex)
	fm.fileCachedMemoryBytes = make(map[IFile][]bool)
	fm.allFiles = make([]IFile, 0)

	return nil
}

func (fm *FileMirror) ReadAt(
	b []byte,
	off int64,
	operationUserData any,
	useFiles ...IFile,
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
		operation.fileMirror = fm

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
	useFiles ...IFile,
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
		operation.fileMirror = fm

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

func (fm *FileMirror) GetAsyncFiles() []IFile {
	files := make([]IFile, 0)

	for _, file := range fm.allFiles {
		if fm.asyncFiles[file] {
			files = append(files, file)
		}
	}

	return files
}

func (fm *FileMirror) GetFirstAsyncFile() IFile {
	for file, _ := range fm.asyncFiles {
		return file
	}

	return nil
}

func (fm *FileMirror) GetNonAsyncFiles() []IFile {
	files := make([]IFile, 0)

	for _, file := range fm.allFiles {
		if !fm.asyncFiles[file] {
			files = append(files, file)
		}
	}

	return files
}

func (fm *FileMirror) GetFirstNonAsyncFile() IFile {
	for _, file := range fm.allFiles {
		if !fm.asyncFiles[file] {
			return file
		}
	}

	return nil
}

func (fm *FileMirror) GetAllFiles() []IFile {
	return fm.allFiles
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

func (fm *FileMirror) HasAsyncOperations() bool {
	return len(fm.asyncOperations) > 0
}

func (fm *FileMirror) GetCountAsyncOperations() int {
	return len(fm.asyncOperations)
}

func NewFileMirror(queueSize int, idleSleepDuration ...time.Duration) *FileMirror {
	fm := FileMirror{}
	fm.asyncOperations = make(chan *Operation, queueSize)
	fm.fileMutexes = make(map[IFile]*sync.Mutex)
	fm.asyncFiles = make(map[IFile]bool)
	fm.fileUserData = make(map[IFile]any)
	fm.fileCachedMemoryBytes = make(map[IFile][]bool)

	if len(idleSleepDuration) == 0 {
		fm.idleSleepDuration = 10 * time.Millisecond
	} else {
		fm.idleSleepDuration = idleSleepDuration[0]
	}

	go fm.run()

	return &fm
}
