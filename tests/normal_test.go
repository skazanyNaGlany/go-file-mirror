package tests

import (
	"crypto/rand"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestNormal(t *testing.T) {
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

	operationCallbackCount := 0
	idleCallbackCount := 0

	buffer := make([]byte, 10)

	fm.SetOperationCallback(func(operation *gofilemirror.Operation) {
		switch operationCallbackCount {
		case 0:
			// started reading from 0 position
			// no data available since the file is empty
			assert.Equal(t, gofilemirror.OT_READ_AT, operation.GetType())
			assert.False(t, operation.IsAsync())
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
			assert.False(t, operation.IsAsync())
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
			// started writing to 0 position (file "f")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.False(t, operation.IsAsync())
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
			// done writing to 0 position (file "f")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.False(t, operation.IsAsync())
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
			// started writing to 0 position (file "f2")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.False(t, operation.IsAsync())
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
			// done writing to 0 position (file "f2")
			assert.Equal(t, gofilemirror.OT_WRITE_AT, operation.GetType())
			assert.False(t, operation.IsAsync())
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
			assert.False(t, operation.IsAsync())
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
			assert.False(t, operation.IsAsync())
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

		operationCallbackCount++
	})

	fm.SetIdleCallback(func(fileMirror *gofilemirror.FileMirror) {
		idleCallbackCount++
	})

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	fm.SetFileMutex(f, &sync.Mutex{})
	fm.SetFileMutex(f2, &sync.Mutex{})

	fm.SetFileUserData(f, "some_user_data_here1")
	fm.SetFileUserData(f2, "some_user_data_here2")

	assert.Contains(t, fm.GetNonAsyncFiles(), f)
	assert.Contains(t, fm.GetNonAsyncFiles(), f2)
	assert.NotContains(t, fm.GetAsyncFiles(), f)
	assert.NotContains(t, fm.GetAsyncFiles(), f2)

	fm.SetFileCachedMemoryBytes(
		f,
		make([]bool, len(buffer)))

	fm.SetFileCachedMemoryBytes(
		f2,
		make([]bool, len(buffer)))

	assert.Len(t, fm.GetFileCachedMemoryBytes(f), 10)
	assert.Len(t, fm.GetFileCachedMemoryBytes(f2), 10)

	assert.False(t, fm.IsFileFullyCached(f))
	assert.False(t, fm.IsFileFullyCached(f2))

	assert.Equal(t, 0, fm.GetFileCachedPercent(f))
	assert.Equal(t, 0, fm.GetFileCachedPercent(f2))

	// read from only one file "f"
	// the following operation, as well as all
	// operations in this test will block untill
	// all data is readed from all files
	// (currently we are reading from one file, writing to
	// two files)
	operations := fm.ReadAt(buffer, 0, "some_user_data_here")

	assert.Equal(t, 2, operationCallbackCount)
	assert.NotNil(t, operations.GetNonAsyncOperations())
	assert.NotNil(t, operations.GetAsyncOperations())
	assert.Len(t, *operations.GetNonAsyncOperations(), 1)
	assert.Len(t, *operations.GetAsyncOperations(), 0)
	assert.Empty(t, operations.GetPendingOperations())

	// this is not necessary since we are reading the data
	// in sync mode, so all operations should done immediately
	// I left it here for informative
	operations.GetNonAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForDone(1 * time.Second)

	// no data to read since the file is empty
	firstOperation := (*operations.GetNonAsyncOperations())[0]
	assert.Zero(t, firstOperation.GetLastResultInt())
	assert.ErrorAs(t, firstOperation.GetLastResultError(), &io.EOF)
	assert.True(t, firstOperation.IsRead())

	// fill the buffer with some random data
	rand.Read(buffer)

	operations = fm.WriteAt(buffer, 0, "some_user_data_here")
	assert.Equal(t, 6, operationCallbackCount)
	assert.NotNil(t, operations.GetNonAsyncOperations())
	assert.NotNil(t, operations.GetAsyncOperations())
	assert.Len(t, *operations.GetNonAsyncOperations(), 2)
	assert.Len(t, *operations.GetAsyncOperations(), 0)
	assert.Empty(t, operations.GetPendingOperations())

	operations.GetNonAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForDone(1 * time.Second)

	firstOperation = (*operations.GetNonAsyncOperations())[0]
	assert.True(t, firstOperation.IsWrite())

	assert.False(t, fm.IsFileFullyCached(f))
	assert.False(t, fm.IsFileFullyCached(f2))

	// read written data
	operations = fm.ReadAt(buffer, 0, "some_user_data_here")

	assert.Equal(t, 8, operationCallbackCount)
	assert.NotNil(t, operations.GetNonAsyncOperations())
	assert.NotNil(t, operations.GetAsyncOperations())
	assert.Len(t, *operations.GetNonAsyncOperations(), 1)
	assert.Len(t, *operations.GetAsyncOperations(), 0)
	assert.Empty(t, operations.GetPendingOperations())

	operations.GetNonAsyncOperations().WaitForStart(1 * time.Second)
	operations.GetNonAsyncOperations().WaitForDone(1 * time.Second)

	// no data to read since the file is empty
	firstOperation = (*operations.GetNonAsyncOperations())[0]
	assert.Equal(t, int64(10), firstOperation.GetLastResultInt())
	assert.Nil(t, firstOperation.GetLastResultError())
	assert.True(t, firstOperation.IsRead())

	assert.True(t, fm.IsFileFullyCached(f))
	assert.False(t, fm.IsFileFullyCached(f2))
	assert.NotZero(t, idleCallbackCount)

	assert.Equal(t, 100, fm.GetFileCachedPercent(f))
}
