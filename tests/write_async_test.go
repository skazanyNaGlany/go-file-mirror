package tests

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestWriteAsync(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

	f, err := os.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := os.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	fm.SetReadingFile(f)
	fm.SetFileAsync(f2, true)
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	fm.SetFileMutex(f, &sync.Mutex{})
	fm.SetFileMutex(f2, &sync.Mutex{})

	callbackCalledCount := 0

	fm.SetAsyncOperationCallback(func(operation *gofilemirror.AsyncOperation) bool {
		switch callbackCalledCount {
		case 0:
			assert.Equal(t, gofilemirror.AOT_WRITE_STRING, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.NotEmpty(t, operation.GetStringBuffer())
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
		case 1:
			assert.Equal(t, gofilemirror.AOT_WRITE_STRING, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.NotEmpty(t, operation.GetStringBuffer())
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
		case 2:
			assert.Equal(t, gofilemirror.AOT_SYNC, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
		case 3:
			assert.Equal(t, gofilemirror.AOT_SYNC, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
		case 4:
			assert.Equal(t, gofilemirror.AOT_SEEK, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.Equal(t, io.SeekStart, operation.GetWhence())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
		case 5:
			assert.Equal(t, gofilemirror.AOT_SEEK, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.Equal(t, io.SeekStart, operation.GetWhence())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
		case 6:
			assert.Equal(t, gofilemirror.AOT_TRUNCATE, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.Equal(t, int64(2), operation.GetSize())
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
		case 7:
			assert.Equal(t, gofilemirror.AOT_TRUNCATE, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.Equal(t, int64(2), operation.GetSize())
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
		case 8:
			assert.Equal(t, gofilemirror.AOT_WRITE_AT, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.NotNil(t, operation.GetBuffer())
			assert.Equal(t, int64(2), operation.GetOffset())
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
		case 9:
			assert.Equal(t, gofilemirror.AOT_WRITE_AT, operation.GetType())
			assert.Equal(t, f2, operation.GetFile())
			assert.NotNil(t, operation.GetBuffer())
			assert.Equal(t, int64(2), operation.GetOffset())
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
		}

		callbackCalledCount++

		return true
	})

	assert.NotNil(t, fm.GetAsyncOperationCallback())

	strb := "test123"
	strb2 := []byte("789def")

	// case 0-1
	ops, n, err := fm.WriteString(strb, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 2, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	assert.Nil(t, ops[0].GetLastResultError())
	assert.Equal(t, int(ops[0].GetLastResultInt()), 7)
	assert.Equal(t, ops[0].GetStringBuffer(), strb)

	// case 2-3
	ops, err = fm.Sync(nil)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 4, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	// set file position to 0
	// case 4-5
	ops, n2, err := fm.Seek(0, io.SeekStart, nil)

	assert.Zero(t, n2)
	assert.Nil(t, err)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 6, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	// read at 0 position
	readed := make([]byte, len(strb))

	ops, n, err = fm.Read(readed, nil)
	assert.Nil(t, err)
	assert.Equal(t, n, len(strb))
	assert.Empty(t, ops)
	assert.Equal(t, strb, string(readed))

	// case 6-7
	ops, err = fm.Truncate(2, nil)
	assert.Nil(t, err)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 8, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	// case 8-9
	ops, n, err = fm.WriteAt(strb2, 2, nil)
	assert.Nil(t, err)
	assert.Len(t, ops, 1)
	assert.Equal(t, len(strb2), n)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 10, callbackCalledCount)
	assert.True(t, ops[0].IsDone())
}
