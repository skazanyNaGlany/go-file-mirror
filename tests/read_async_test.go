package tests

import (
	"io"
	"os"
	"testing"
	"time"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestReadAsync(t *testing.T) {
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
	assert.True(t, fm.AddAsyncFile(f))

	callbackCalledCount := 0

	fm.SetAsyncOperationCallback(func(operation *gofilemirror.AsyncOperation) {
		callbackCalledCount++
	})

	// read at 0 position
	readed := make([]byte, 6)

	ops, n, err := f.Read(readed)
	assert.Nil(t, err)
	assert.Zero(t, n)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 2, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	// no data to read from empty file
	assert.ErrorAs(t, ops[0].GetLastResultError(), &io.EOF)
	assert.Zero(t, ops[0].GetLastResultInt())
	assert.Equal(t, ops[0].GetBuffer(), make([]byte, 6))

	// set file position to 0
	n2, err := f.Seek(0, io.SeekStart)

	assert.Zero(t, n2)
	assert.Nil(t, err)

	// write some test data
	strb := []byte("123abc")

	n, err = f.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)

	// set file position to 0
	n2, err = f.Seek(0, io.SeekStart)

	assert.Zero(t, n2)
	assert.Nil(t, err)

	// read again, this time with data in the file
	ops, n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Zero(t, n)
	assert.Len(t, ops, 1)

	ops[0].WaitForStart(10 * time.Second)
	assert.True(t, ops[0].IsStarted())

	ops[0].WaitForDone(10 * time.Second)
	assert.Equal(t, 4, callbackCalledCount)
	assert.True(t, ops[0].IsDone())

	// no data to read from empty file
	assert.Nil(t, err)
	assert.Equal(t, len(strb), int(ops[0].GetLastResultInt()))
	assert.Equal(t, ops[0].GetBuffer(), strb)

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
