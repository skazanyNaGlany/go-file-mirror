package tests

import (
	"os"
	"sync"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
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

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	assert.Len(t, fm.GetAllFiles(), 2)
	assert.Contains(t, fm.GetAllFiles(), f)
	assert.Contains(t, fm.GetAllFiles(), f2)

	assert.Contains(t, fm.GetReadingFiles(), f)
	assert.Contains(t, fm.GetWritingFiles(), f)
	assert.Contains(t, fm.GetWritingFiles(), f2)

	fm.SetFileAsync(f, true)
	fm.SetFileAsync(f2, true)

	assert.True(t, fm.IsFileAsync(f))
	assert.True(t, fm.IsFileAsync(f2))

	assert.Nil(t, fm.GetFileMutex(f))
	assert.Nil(t, fm.GetFileMutex(f2))

	fm.SetFileMutex(f, &sync.Mutex{})
	fm.SetFileMutex(f2, &sync.Mutex{})

	assert.NotNil(t, fm.GetFileMutex(f))
	assert.NotNil(t, fm.GetFileMutex(f2))

	assert.True(t, fm.RemoveReadingFile(f))
	assert.True(t, fm.RemoveWritingFile(f))
	assert.True(t, fm.RemoveWritingFile(f2))
	fm.SetFileAsync(f, false)
	fm.SetFileAsync(f2, false)

	assert.NotContains(t, fm.GetReadingFiles(), f)
	assert.NotContains(t, fm.GetWritingFiles(), f)
	assert.NotContains(t, fm.GetWritingFiles(), f2)
	assert.False(t, fm.IsFileAsync(f))
	assert.False(t, fm.IsFileAsync(f2))

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	fm.RemoveAllFiles()
	assert.Empty(t, fm.GetAllFiles())

	// all files needs to be addes to the FileMirror instance
	// to close it automatically when closing that FileMirror
	// instance
	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))
}
