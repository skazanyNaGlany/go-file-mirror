package tests

import (
	"io"
	"os"
	"testing"
	"time"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestWriteAsync(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

	f, err := gofilemirror.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := gofilemirror.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddReadingFile(f2))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))
	assert.True(t, fm.AddAsyncFile(f2))

	callbackCalledCount := 0

	fm.SetAsyncOperationCallback(func(operation *gofilemirror.AsyncOperation) {
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
	})

	strb := "test123"
	strb2 := []byte("789def")

	// case 0-1
	ops, n, err := f2.WriteString(strb)
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
	ops, err = f2.Sync()
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
	ops, n2, err := f.Seek(0, io.SeekStart)

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

	ops, n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, len(strb))
	assert.Empty(t, ops)
	assert.Equal(t, strb, string(readed))

	// case 6-7
	ops, err = f2.Truncate(2)
	assert.Nil(t, err)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 8, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	// case 8-9
	ops, n, err = f.WriteAt(strb2, 2)
	assert.Nil(t, err)
	assert.Len(t, ops, 1)
	assert.Equal(t, len(strb2), n)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 10, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
