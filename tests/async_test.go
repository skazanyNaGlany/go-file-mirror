package tests

import (
	"crypto/rand"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestAsync(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close(true)

	f, err := os.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := os.CreateTemp("/tmp", "testing_file_mirror2")
	if err != nil {
		panic(err)
	}

	callbackCalledCount := 0
	idleCallbackCalledCount := 0

	buffer := make([]byte, 10)

	fm.SetIdleCallback(func(fileMirror *gofilemirror.FileMirror) {
		idleCallbackCalledCount++
	})

	fm.SetOperationCallback(func(operation *gofilemirror.Operation) {
		switch callbackCalledCount {
		case 0:
			// started reading from 0 position
			// no data available since the file is empty
			assert.Equal(t, gofilemirror.OT_READ_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f, operation.GetFile())
			assert.NotNil(t, f, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here1", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Zero(t, operation.GetLastResultInt())
		case 1:
			// done reading from 0 position
			// no data available since the file is empty
			// err = io.EOF
			assert.Equal(t, gofilemirror.OT_READ_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f, operation.GetFile())
			assert.NotNil(t, f, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here1", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.ErrorAs(t, operation.GetLastResultError(), &io.EOF)
			assert.Zero(t, operation.GetLastResultInt())
		case 2:
			// started writing to 0 position (file "f" and "f2")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f, operation.GetFile())
			assert.NotNil(t, f, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here1", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Zero(t, operation.GetLastResultInt())
		case 3:
			// done writing to 0 position (file "f" and "f2")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f, operation.GetFile())
			assert.NotNil(t, f, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here1", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Equal(t, int64(len(buffer)), operation.GetLastResultInt())
		case 4:
			// started writing to 0 position (file "f" and "f2")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f2, operation.GetFile())
			assert.NotNil(t, f2, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here2", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Zero(t, operation.GetLastResultInt())
		case 5:
			// done writing to 0 position (file "f" and "f2")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f2, operation.GetFile())
			assert.NotNil(t, f2, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here2", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Equal(t, int64(len(buffer)), operation.GetLastResultInt())
		case 6:
			// started reading from 0 position
			assert.Equal(t, gofilemirror.OT_READ_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f, operation.GetFile())
			assert.NotNil(t, f, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.False(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here1", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Zero(t, operation.GetLastResultInt())
		case 7:
			// done reading from 0 position
			assert.Equal(t, gofilemirror.OT_READ_AT, operation.GetType())
			assert.True(t, operation.IsAsync())
			assert.Equal(t, f, operation.GetFile())
			assert.NotNil(t, f, operation.GetBuffer())
			assert.Equal(t, int64(0), operation.GetOffset())
			assert.Len(t, operation.GetBuffer(), len(buffer))
			assert.True(t, operation.IsStarted())
			assert.True(t, operation.IsDone())
			assert.Equal(t, "some_user_data_here1", operation.GetFileUserData())
			assert.Equal(t, "some_user_data_here", operation.GetUserData())
			assert.Nil(t, operation.GetLastResultError())
			assert.Equal(t, int64(len(buffer)), operation.GetLastResultInt())
		}

		operation.WaitForStart(0 * time.Second)
		operation.WaitForDone(0 * time.Second)

		callbackCalledCount++
	})

	// reading only from "f"
	// writing to "f" and "f2"
	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	fm.SetFileMutex(f, &sync.Mutex{})
	fm.SetFileMutex(f2, &sync.Mutex{})

	// all operations on both files will be async
	fm.SetFileAsync(f, true)
	fm.SetFileAsync(f2, true)

	assert.Equal(t, f, fm.GetFirstAsyncFile())
	assert.Nil(t, fm.GetFirstNonAsyncFile())

	assert.Contains(t, fm.GetAsyncFiles(), f)
	assert.Contains(t, fm.GetAsyncFiles(), f2)
	assert.NotContains(t, fm.GetNonAsyncFiles(), f)
	assert.NotContains(t, fm.GetNonAsyncFiles(), f2)

	fm.SetFileUserData(f, "some_user_data_here1")
	fm.SetFileUserData(f2, "some_user_data_here2")

	// read from only one file "f"
	// the following operation, as well as all
	// operations in this test will non-block
	operations := fm.ReadAt(buffer, 0, "some_user_data_here")

	time.Sleep(1 * time.Second)
	fm.WaitForNoAsyncOperations(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForDone(1 * time.Second)
	operations.GetAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetAsyncOperations().WaitForDone(1 * time.Second)

	// no data to read since the file is empty
	firstOperation := (*operations.GetAsyncOperations())[0]
	assert.Zero(t, firstOperation.GetLastResultInt())
	assert.ErrorAs(t, firstOperation.GetLastResultError(), &io.EOF)

	assert.Equal(t, 2, callbackCalledCount)
	assert.NotNil(t, operations.GetNonAsyncOperations())
	assert.NotNil(t, operations.GetAsyncOperations())
	assert.Len(t, *operations.GetNonAsyncOperations(), 0)
	assert.Len(t, *operations.GetAsyncOperations(), 1)
	assert.Empty(t, operations.GetPendingOperations())

	// fill the buffer with some random data
	rand.Read(buffer)

	operations = fm.WriteAt(buffer, 0, "some_user_data_here")

	time.Sleep(1 * time.Second)
	fm.WaitForNoAsyncOperations(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForDone(1 * time.Second)
	operations.GetAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetAsyncOperations().WaitForDone(1 * time.Second)

	assert.Equal(t, 6, callbackCalledCount)
	assert.NotNil(t, operations.GetNonAsyncOperations())
	assert.NotNil(t, operations.GetAsyncOperations())
	assert.Len(t, *operations.GetNonAsyncOperations(), 0)
	assert.Len(t, *operations.GetAsyncOperations(), 2)
	assert.Empty(t, operations.GetPendingOperations())

	firstOperation = (*operations.GetAsyncOperations())[0]
	assert.Equal(t, int64(10), firstOperation.GetLastResultInt())
	assert.Nil(t, firstOperation.GetLastResultError())

	secondOperation := (*operations.GetAsyncOperations())[1]
	assert.Equal(t, int64(10), secondOperation.GetLastResultInt())
	assert.Nil(t, secondOperation.GetLastResultError())

	// read written data
	operations = fm.ReadAt(buffer, 0, "some_user_data_here")

	time.Sleep(1 * time.Second)
	fm.WaitForNoAsyncOperations(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForDone(1 * time.Second)
	operations.GetAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetAsyncOperations().WaitForDone(1 * time.Second)

	assert.Equal(t, 8, callbackCalledCount)
	assert.NotNil(t, operations.GetNonAsyncOperations())
	assert.NotNil(t, operations.GetAsyncOperations())
	assert.Len(t, *operations.GetNonAsyncOperations(), 0)
	assert.Len(t, *operations.GetAsyncOperations(), 1)
	assert.Empty(t, operations.GetPendingOperations())

	firstOperation = (*operations.GetAsyncOperations())[0]
	assert.Equal(t, int64(10), firstOperation.GetLastResultInt())
	assert.Nil(t, firstOperation.GetLastResultError())

	assert.NotZero(t, idleCallbackCalledCount)

	log.Println(idleCallbackCalledCount)
}
