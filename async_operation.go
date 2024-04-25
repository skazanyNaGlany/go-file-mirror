package gofilemirror

import (
	"os"
	"time"
)

type AsyncOperation struct {
	file         *os.File
	_type        AsyncOperationType
	err          error
	resultInt    int64
	buffer       []byte
	offset       int64
	whence       int
	size         int64
	stringBuffer string
	started      bool
	done         bool
}

func (ao *AsyncOperation) GetFile() *os.File {
	return ao.file
}

func (ao *AsyncOperation) GetType() AsyncOperationType {
	return ao._type
}

func (ao *AsyncOperation) GetLastResultError() error {
	return ao.err
}

func (ao *AsyncOperation) GetLastResultInt() int64 {
	return ao.resultInt
}

func (ao *AsyncOperation) GetBuffer() []byte {
	return ao.buffer
}

func (ao *AsyncOperation) GetOffset() int64 {
	return ao.offset
}

func (ao *AsyncOperation) GetWhence() int {
	return ao.whence
}

func (ao *AsyncOperation) GetSize() int64 {
	return ao.size
}

func (ao *AsyncOperation) GetStringBuffer() string {
	return ao.stringBuffer
}

func (ao *AsyncOperation) IsStarted() bool {
	return ao.started
}

func (ao *AsyncOperation) IsDone() bool {
	return ao.done
}

func (ao *AsyncOperation) WaitForStart(duration time.Duration) {
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

func (ao *AsyncOperation) WaitForDone(duration time.Duration) {
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
