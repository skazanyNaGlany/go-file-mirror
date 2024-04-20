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
		callbackCalledCount++
	})

	strb := "test123"
	strb2 := []byte("789def")

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

	ops, err = f2.Truncate(2)
	assert.Nil(t, err)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 8, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

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
