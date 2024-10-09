package gofilemirror

import (
	"time"
)

type Operation struct {
	file              IFile
	_type             OperationType
	err               error
	resultInt         int64
	buffer            []byte
	offset            int64
	started           bool
	done              bool
	async             bool
	operationUserData any
	fileUserData      any
	fileMirror        *FileMirror
	timeMilisecond    int64
	// for Seek, Stat and WriteString
	// whence            int
	// size              int64
	// stringBuffer      string
}

func (ao *Operation) GetTimeMilisecond() int64 {
	return ao.timeMilisecond
}

func (ao *Operation) SetTimeMilisecond(timeMilisecond int64) {
	ao.timeMilisecond = timeMilisecond
}

func (ao *Operation) GetFileMirror() *FileMirror {
	return ao.fileMirror
}

func (ao *Operation) SetFileMirror(fileMirror *FileMirror) {
	ao.fileMirror = fileMirror
}

func (ao *Operation) GetUserData() any {
	return ao.operationUserData
}

func (ao *Operation) SetUserData(userData any) {
	ao.operationUserData = userData
}

func (ao *Operation) GetFileUserData() any {
	return ao.fileUserData
}

func (ao *Operation) SetFileUserData(userData any) {
	ao.fileUserData = userData
}

func (ao *Operation) GetFile() IFile {
	return ao.file
}

func (ao *Operation) SetFile(file IFile) {
	ao.file = file
}

func (ao *Operation) GetType() OperationType {
	return ao._type
}

func (ao *Operation) SetType(_type OperationType) {
	ao._type = _type
}

func (ao *Operation) GetLastResultError() error {
	return ao.err
}

func (ao *Operation) SetLastResultError(err error) {
	ao.err = err
}

func (ao *Operation) GetLastResultInt() int64 {
	return ao.resultInt
}

func (ao *Operation) SetLastResultInt(resultInt int64) {
	ao.resultInt = resultInt
}

func (ao *Operation) GetBuffer() []byte {
	return ao.buffer
}

func (ao *Operation) SetBuffer(buffer []byte) {
	ao.buffer = buffer
}

func (ao *Operation) GetOffset() int64 {
	return ao.offset
}

func (ao *Operation) SetOffset(offset int64) {
	ao.offset = offset
}

func (ao *Operation) IsStarted() bool {
	return ao.started
}

func (ao *Operation) SetStarted(started bool) {
	ao.started = started
}

func (ao *Operation) SetAsync(async bool) {
	ao.async = async
}

func (ao *Operation) IsDone() bool {
	return ao.done
}

func (ao *Operation) IsAsync() bool {
	return ao.async
}

func (ao *Operation) IsRead() bool {
	return ao._type == OT_READ_AT
}

func (ao *Operation) IsWrite() bool {
	return ao._type == OT_WRITE_AT
}

func (ao *Operation) SetDone(started bool) {
	ao.started = started
}

func (ao *Operation) WaitForStart(duration time.Duration) {
	currentDuration := time.Duration(0)

	for {
		if ao.started {
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

func (ao *Operation) WaitForDone(duration time.Duration) {
	currentDuration := time.Duration(0)

	for {
		if ao.done {
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

// func (ao *Operation) GetWhence() int {
// 	return ao.whence
// }

// func (ao *Operation) SetWhence(whence int) {
// 	ao.whence = whence
// }

// func (ao *Operation) GetSize() int64 {
// 	return ao.size
// }

// func (ao *Operation) SetSize(size int64) {
// 	ao.size = size
// }

// func (ao *Operation) GetStringBuffer() string {
// 	return ao.stringBuffer
// }

// func (ao *Operation) SetStringBuffer(stringBuffer string) {
// 	ao.stringBuffer = stringBuffer
// }
