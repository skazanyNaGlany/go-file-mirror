package tests

import (
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
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
	assert.True(t, fm.AddAsyncFile(f2))

	assert.Contains(t, fm.GetReadingFiles(), f)
	assert.Contains(t, fm.GetReadingFiles(), f2)
	assert.Contains(t, fm.GetWritingFiles(), f)
	assert.Contains(t, fm.GetWritingFiles(), f2)
	assert.Contains(t, fm.GetAsyncFiles(), f)
	assert.Contains(t, fm.GetAsyncFiles(), f2)

	assert.True(t, fm.RemoveReadingFile(f))
	assert.True(t, fm.RemoveReadingFile(f2))
	assert.True(t, fm.RemoveWritingFile(f))
	assert.True(t, fm.RemoveWritingFile(f2))
	assert.True(t, fm.RemoveAsyncFile(f))
	assert.True(t, fm.RemoveAsyncFile(f2))

	assert.NotContains(t, fm.GetReadingFiles(), f)
	assert.NotContains(t, fm.GetReadingFiles(), f2)
	assert.NotContains(t, fm.GetWritingFiles(), f)
	assert.NotContains(t, fm.GetWritingFiles(), f2)
	assert.NotContains(t, fm.GetAsyncFiles(), f)
	assert.NotContains(t, fm.GetAsyncFiles(), f2)

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddReadingFile(f2))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))
	assert.True(t, fm.AddAsyncFile(f))
	assert.True(t, fm.AddAsyncFile(f2))

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}