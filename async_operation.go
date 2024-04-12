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

func (ao *AsyncOperation) setType(t AsyncOperationType) {
	ao._type = t
}

func (ao *AsyncOperation) setLastResultError(err error) {
	ao.err = err
}

func (ao *AsyncOperation) setLastResultInt(i int64) {
	ao.resultInt = i
}

func (ao *AsyncOperation) setBuffer(buff []byte) {
	ao.buff = buff
}

func (ao *AsyncOperation) setOffset(off int64) {
	ao.off = off
}

func (ao *AsyncOperation) setWhence(whence int) {
	ao.whence = whence
}

func (ao *AsyncOperation) setSize(size int64) {
	ao.size = size
}

func (ao *AsyncOperation) setStringBuffer(s string) {
	ao.stringBuffer = s
}
