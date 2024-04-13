package gofilemirror

type AsyncOperation struct {
	_type        AsyncOperationType
	err          error
	resultInt    int64
	buff         []byte
	off          int64
	whence       int
	size         int64
	stringBuffer string
	started      bool
	done         bool
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
	return ao.buff
}

func (ao *AsyncOperation) GetOffset() int64 {
	return ao.off
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
